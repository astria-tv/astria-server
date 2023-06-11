package serve

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // For Profiling
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/goava/di"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/olaris/olaris-server/ffmpeg"
	"gitlab.com/olaris/olaris-server/interfaces/web"
	"gitlab.com/olaris/olaris-server/metadata"
	"gitlab.com/olaris/olaris-server/metadata/agents"
	"gitlab.com/olaris/olaris-server/metadata/app"
	"gitlab.com/olaris/olaris-server/metadata/db"
	"gitlab.com/olaris/olaris-server/react"
	"gitlab.com/olaris/olaris-server/streaming"
)

func Options() di.Option {
	return di.Options(
		di.Provide(NewServeCommand, di.Tags{"type": "serve"}),
		di.Invoke(RegisterServeCommand),
	)
}

func RegisterServeCommand(deps struct {
	di.Inject
	RootCommand  *cobra.Command `di:"type=root"`
	ServeCommand *cobra.Command `di:"type=serve"`
}) {
	deps.RootCommand.AddCommand(deps.ServeCommand)
}

func NewServeCommand(deps struct {
	di.Inject
	StreamingController web.Controller `di:"type=streaming"`
}) *cobra.Command {
	registerMetadataFlags, registerMetadataViper := metadata.FlagFuncs()

	c := &cobra.Command{
		Use:   "serve",
		Short: "Start the olaris server",
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("server.port", cmd.Flags().Lookup("port"))
			_ = viper.BindPFlag("server.dbLog", cmd.Flags().Lookup("db-log"))
			_ = registerMetadataViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetBool("server.verbose") {
				log.SetLevel(log.DebugLevel)
			}

			// Check FFmpeg version and warn if it's missing
			ffmpegVersion, err := ffmpeg.GetFfmpegVersion()
			if err != nil {
				if parseErr, ok := err.(*ffmpeg.VersionParseError); ok {
					log.WithError(parseErr).Warn("unable to determine installed FFmpeg version")
				} else {
					log.WithError(err).Warn("FFmpeg not found. STREAMING WILL NOT WORK IF FFMPEG IS NOT INSTALLED AND IN YOUR PATH!")
				}
			} else {
				log.WithField("version", ffmpegVersion.ToString()).Debugf("FFmpeg found")
			}

			// Check FFprobe version and warn if it's missing
			ffprobeVersion, err := ffmpeg.GetFfprobeVersion()
			if err != nil {
				if parseErr, ok := err.(*ffmpeg.VersionParseError); ok {
					log.WithError(parseErr).Warn("unable to determine installed FFprobe version")
				} else {
					log.WithError(err).Warn("FFprobe not found. STREAMING WILL NOT WORK IF FFPROBE IS NOT INSTALLED AND IN YOUR PATH!")
				}
			} else {
				log.WithField("version", ffprobeVersion.ToString()).Debugf("FFprobe found")
			}

			mainRouter := mux.NewRouter()

			r := mainRouter.PathPrefix("/olaris")
			rr := mainRouter.PathPrefix("/olaris")
			rrr := mainRouter.PathPrefix("/olaris")

			dbOptions := db.DatabaseOptions{
				Connection: viper.GetString("database.connection"),
				LogMode:    viper.GetBool("server.dbLog"),
			}

			mctx := app.NewMDContext(dbOptions, agents.NewTmdbAgent())
			if viper.GetBool("server.verbose") {
				log.SetLevel(log.DebugLevel)
			}
			viper.WatchConfig()

			updateConfig := func(in fsnotify.Event) {
				log.Infoln("configuration file change detected")
				if viper.GetBool("server.verbose") {
					log.SetLevel(log.DebugLevel)
				} else {
					log.SetLevel(log.InfoLevel)
				}
				mctx.Db.LogMode(viper.GetBool("server.dbLog"))
			}
			viper.OnConfigChange(updateConfig)

			metaRouter := r.PathPrefix("/m").Subrouter()
			metadata.RegisterRoutes(mctx, metaRouter)

			streamingRouter := rr.PathPrefix("/s").Subrouter()
			deps.StreamingController.RegisterRoutes(streamingRouter)

			// This is just to make sure that no temp files stay behind in case the
			// garbage collection below didn't work properly for some reason.
			// This is also relevant during development because the realize auto-reload
			// tool doesn't properly send SIGTERM.
			ffmpeg.CleanTranscodingCache()
			port := viper.GetInt("server.port")

			appRoute := rrr.PathPrefix("/app").
				Handler(http.StripPrefix("/olaris/app", react.GetHandler())).
				Name("app")

			appURL, _ := appRoute.URL()
			mainRouter.Path("/").Handler(http.RedirectHandler(appURL.Path, http.StatusMovedPermanently))
			mainRouter.Path("/olaris").Handler(http.RedirectHandler(appURL.Path, http.StatusMovedPermanently))

			handler := cors.AllowAll().Handler(mainRouter)
			handler = handlers.LoggingHandler(os.Stdout, handler)

			log.Infoln("binding on port", port)
			srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					log.WithFields(log.Fields{"error": err}).Fatal("error starting server")
				}
			}()

			stopChan := make(chan os.Signal, 2)
			signal.Notify(stopChan, os.Interrupt, os.Kill)

			// Wait for termination signal
			<-stopChan
			log.Println("shutting down...")

			// Clean up the metadata context
			mctx.Cleanup()

			// Clean up playback/transcode sessions
			sessionCleanupContext, sessionCleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer sessionCleanupCancel()
			streaming.PBSManager.DestroyAll(sessionCleanupContext)

			// Shut down the HTTP server
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)

			log.Println("shut down complete, exiting.")
			return nil
		},
	}

	c.Flags().IntP("port", "p", 8080, "http port")
	c.Flags().Bool("db-log", false, "sets whether the database should log queries")
	c.Flags().Bool("zeroconf-enabled", false, "enables the zeroconf service")
	c.Flags().String("zeroconf-domain", "local.", "sets the domain for the zeroconf service if it is enabled")
	_ = registerMetadataFlags(c)

	return c
}
