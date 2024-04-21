package imagestore

import (
	"fmt"
	"io"
	"os"

	"github.com/matmazurk/acc2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type store struct {
	basepath string
	logger   zerolog.Logger
}

func NewStore(basepath string, logger zerolog.Logger) (store, error) {
	err := os.MkdirAll(basepath+photosRelativeDir, 0o750)
	if err != nil && !os.IsExist(err) {
		return store{}, errors.Wrap(err, "could not create photos dir")
	}
	return store{
		basepath: basepath,
		logger:   logger,
	}, nil
}

const (
	photosRelativeDir  = "/photos"
	filenameTimeLayout = "020106_1504"
)

func (s store) SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error {
	defer r.Close()

	photoPath := s.providePhotoAbsolutePath(e, fileExtension)
	file, err := os.Create(photoPath)
	if err != nil {
		return errors.Wrapf(err, "could not create file '%s'", photoPath)
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return errors.Wrap(err, "could not copy file contents")
	}

	return nil
}

func (s store) dirAbsolutePath() string {
	return s.basepath + photosRelativeDir
}

func (s store) providePhotoAbsolutePath(e model.Expense, fileExtension string) string {
	return fmt.Sprintf("%s/%s_%s.%s", s.dirAbsolutePath(), e.CreatedAt().Format(filenameTimeLayout), e.ID(), fileExtension)
}
