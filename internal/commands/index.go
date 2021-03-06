package commands

import (
	"context"
	"time"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/urfave/cli"
)

// IndexCommand is used to register the index cli command
var IndexCommand = cli.Command{
	Name:   "index",
	Usage:  "Indexes media files in originals path",
	Flags:  indexFlags,
	Action: indexAction,
}

var indexFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all, a",
		Usage: "re-index all originals, including unchanged files",
	},
}

// indexAction indexes all photos in originals directory (photo library)
func indexAction(ctx *cli.Context) error {
	start := time.Now()

	conf := config.NewConfig(ctx)
	service.SetConfig(conf)

	if err := conf.CreateDirectories(); err != nil {
		return err
	}

	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := conf.Init(cctx); err != nil {
		return err
	}

	conf.InitDb()
	log.Infof("indexing photos in %s", conf.OriginalsPath())

	if conf.ReadOnly() {
		log.Infof("read-only mode enabled")
	}

	ind := service.Index()

	var opt photoprism.IndexOptions

	if ctx.Bool("all") {
		opt = photoprism.IndexOptionsAll()
	} else {
		opt = photoprism.IndexOptionsNone()
	}

	files := ind.Start(opt)
	elapsed := time.Since(start)

	log.Infof("indexed %d files in %s", len(files), elapsed)

	conf.Shutdown()

	return nil
}
