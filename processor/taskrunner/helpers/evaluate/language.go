package evaluate

import (
	"github.com/PaesslerAG/gval"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	TaskInput    interface{}
	EnvVars      map[string]interface{}
	PipelineUuid uuid.UUID
}

type noIdentifierPresented struct{}

var lng = gval.NewLanguage(
	gval.Bitmask(),
	gval.Text(),
	gval.PropositionalLogic(),
	gval.JSON(),
	arithmetic,
	text,
	funcs,
	dates,
	varselector,
)

func recoverGvalFunc(fName string) {
	if r := recover(); r != nil {
		log.Errorf("Panic while executing a gval function\" %s\", recovered and skipped\n", fName)
		log.Error(r)
	}
}

func Language() gval.Language {
	return lng
}
