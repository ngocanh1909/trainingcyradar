package main
import (
	"net/http"
	"strconv"
	"github.com/labstack/echo"
	"encoding/json"
	"fmt"
	"time"
)
type WorkerType int
const (
	Unknown = iota
	IronWorker
	MechanicalWorker
	OilWorker
)
func(wt WorkerType) String() string {
	names := [4]string{
		"unknown",
		"iron_worker",
		"mechanical_worker",
		"oil_worker",
	}
	if wt < Unknown || wt > OilWorker {
		return "unknown"
	}
	return names[wt]
}
func(wt *WorkerType) ParseFrom(src string) {
	nameWorkerTypeMap := map[string]WorkerType{
		"unknown": Unknown,
		"iron_worker": IronWorker,
		"mechanical_worker": MechanicalWorker,
		"oil_worker": OilWorker,
	}
	if val, exist := nameWorkerTypeMap[src]; exist {
		*wt = val
		return
	}
	*wt = Unknown
}
type Worker struct {
	ID   int              `json:"id"`
	Name string           `json:"name"`
	Age  int              `json:"age"`
	Type WorkerType       `json:"worker_type"`
	OnBoardDate time.Time `json:"-"`
}
type Alias Worker
type AuxWorker struct {
	Type        string     `json:"worker_type"`
	OnBoardDate string     `json:"on_board_date"`
	*Alias
}
func (m *Worker) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&AuxWorker{
			Type: m.Type.String(),
			OnBoardDate: string([]rune(m.OnBoardDate.String())[0:10]),
			Alias: (*Alias)(m),
		},
	)
}
func (m *Worker) UnmarshalJSON(data []byte) error {
	aux := AuxWorker{ Alias: (*Alias)(m)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// WorkerType
	t := new(WorkerType)
	t.ParseFrom(aux.Type)
	m.Type = *t
	// OnBoardDate
	layout := "2006-01-02"
	if t, err := time.Parse(layout, aux.OnBoardDate); err != nil {
		return err
	} else {
		m.OnBoardDate = t
	}
	return nil
}
var (
	Workers = map[int]*Worker{}
	seq   = 1
)
func createWorker(c echo.Context) error {
	u := &Worker{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	Workers[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}
func getWorker(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, Workers[id])
}
func getWorkers(c echo.Context) error {
	fmt.Println(Workers[1])
	return c.JSON(http.StatusOK, Workers)
}
func updateWorker(c echo.Context) error {
	u := new(Worker)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	Workers[id].Name = u.Name
	return c.JSON(http.StatusOK, Workers[id])
}
func deleteWorker(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(Workers, id)
	return c.NoContent(http.StatusNoContent)
}
func main() {
	e := echo.New()
	// Routes
	e.POST("/workers", createWorker)
	e.GET("/workers/:id", getWorker)
	e.GET("/workers", getWorkers)
	e.PUT("/workers/:id", updateWorker)
	e.DELETE("/workers/:id", deleteWorker)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}