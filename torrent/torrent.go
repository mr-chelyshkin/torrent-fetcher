package torrent

import (

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Torrent struct {
	client *torrent.Client
}

// NewTorrent return Torrent object
func NewTorrent() (*Torrent, error) {
	cfg := newConfig()

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Torrent{
		client: client,
	}, nil
}

// DownloadFromFile download torrent data
func (t Torrent) DownloadFromFile(filePath string) error {
	log.Println("start download:", filePath)

	tf, err := t.client.AddTorrentFromFile(filePath)
	if err != nil {
		return errors.Wrap(err, "can't add torrent file to client")
	}

	defer tf.Drop()
	tf.DownloadAll()
	time.Sleep(time.Second*5)

PREPARE:
	for {
		switch tf.PieceStateRuns()[0].Length {
		case tf.NumPieces():
			log.Println(filePath, "already download")
			return nil
		default:
			break PREPARE
		}
	}

WAIT:
	for {
		switch tf.PieceStateRuns()[0].Length {
		case tf.NumPieces():
			log.Println("downloading", filePath, ": complete", tf.BytesCompleted(), "of", tf.Info().TotalLength(), "Bytes")
			log.Println(tf.Info().Name, "downloaded")
			break WAIT
		default:
			log.Println("downloading", filePath, ": complete", tf.BytesCompleted(), "of", tf.Info().TotalLength(), "Bytes")
			time.Sleep(time.Second * 10)
		}
	}

	return nil
}

// gracefully close Torrent object process
func (t Torrent) Close() {
	t.client.Close()
}