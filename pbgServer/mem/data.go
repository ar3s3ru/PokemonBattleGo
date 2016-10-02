package mem

import (
    "github.com/ar3s3ru/PokemonBattleGo/pbgServer"
    "sync"
    "errors"
    "gopkg.in/mgo.v2/bson"
    "time"
    "encoding/json"
    "io/ioutil"
)

type (
    DataBuilder interface {
        UsePokèmonFile(path string)  DataBuilder
        UseTrainersFile(path string) DataBuilder

        Build() pbgServer.IDataMechanism
    }

    memData struct {
        pokèdx   []pbgServer.Pokèmon
        movedx   []pbgServer.Move
        // NOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO!
        // EXTREME BOTTLENECK HERE!
        trainMutex sync.Mutex
        trainers   map[bson.ObjectId]pbgServer.Trainer
    }

    memDataBuilder struct {
        pokèmonFile  string
        trainersFile string
    }
)

var (
    ErrPokèmonNotFound    = errors.New("Pokèmon not found")
    ErrMoveNotFound       = errors.New("Move not found")
    ErrTrainerNotFound    = errors.New("Trainer not found")
    ErrIllegalTrainer     = errors.New("Trainer object is nil")
    ErrInvalidTrainerName = errors.New("Invalid Trainer name used")
)

func NewDataBuilder() DataBuilder {
    return &memDataBuilder{}
}

func (builder *memDataBuilder) UsePokèmonFile(path string) DataBuilder {
    builder.pokèmonFile = path
    return builder
}

func (builder *memDataBuilder) UseTrainersFile(path string) DataBuilder {
    builder.trainersFile = path
    return builder
}

func (builder *memDataBuilder) Build() pbgServer.IDataMechanism {
    if builder.pokèmonFile == "" {
        panic("Using \"\" as Pokèmon file is not allowed")
    } else if file, err := ioutil.ReadFile(builder.pokèmonFile); err != nil {
        panic(err)
    } else {
        pkms := struct {
            Generation int       `json:"generation"`
            PNumbers   int       `json:"pokemon_count"`
            MNumbers   int       `json:"move_count"`
            PList      []pokèmon `json:"pokemons"`
            MList      []move    `json:"moves"`
        }{}

        if err := json.Unmarshal(file, &pkms); err != nil {
            panic(err)
        }

        return &memData{
            pokèdx:   convertLtoPL(pkms.PList),
            movedx:   convertLtoML(pkms.MList),
            trainers: make(map[bson.ObjectId]pbgServer.Trainer),
        }
    }
}

func (data *memData) AddTrainer(trainer pbgServer.Trainer) (bson.ObjectId, error) {
    if trainer == nil {
        return "", ErrIllegalTrainer
    }

    data.trainMutex.Lock()
    defer data.trainMutex.Unlock()

    id := bson.NewObjectIdWithTime(time.Now())
    data.trainers[id] = trainer
    return id, nil
}

func (data *memData) RemoveTrainer(id bson.ObjectId) error {
    data.trainMutex.Lock()
    defer data.trainMutex.Unlock()

    if trainer := data.trainers[id]; trainer == nil {
        return ErrTrainerNotFound
    } else {
        delete(data.trainers, id)
        return nil
    }
}

func (data *memData) GetPokèmons() []pbgServer.Pokèmon {
    return data.pokèdx
}

func (data *memData) GetMoves() []pbgServer.Move {
    return data.movedx
}

func (data *memData) GetMoveById(id int) (pbgServer.Move, error) {
    if id <= 0 || id > len(data.movedx) {
        return nil, ErrMoveNotFound
    } else {
        return data.movedx[id - 1], nil
    }
}

func (data *memData) GetPokèmonById(id int) (pbgServer.Pokèmon, error) {
    if id <= 0 || id > len(data.pokèdx) {
        return nil, ErrPokèmonNotFound
    } else {
        return data.pokèdx[id - 1], nil
    }
}

func (data *memData) GetTrainerById(id bson.ObjectId) (pbgServer.Trainer, error) {
    data.trainMutex.Lock()
    defer data.trainMutex.Unlock()

    if trainer := data.trainers[id]; trainer == nil {
        return nil, ErrTrainerNotFound
    } else {
        return trainer, nil
    }
}

func (data *memData) GetTrainerByName(name string) (pbgServer.Trainer, error) {
    if name == "" {
        return nil, ErrInvalidTrainerName
    }

    data.trainMutex.Lock()
    defer data.trainMutex.Unlock()

    for _, trainer := range data.trainers {
        if trainer.GetName() == name {
            return trainer, nil
        }
    }

    return nil, ErrTrainerNotFound
}
