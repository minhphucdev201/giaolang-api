package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Call struct {
    ID          string    `json:"idcalls"`
    ParentName  string    `json:"parent_name"`
    StudentName string    `json:"student_name"`
    RecordedBy  string    `json:"recorded_by"`
    RecordedAt  time.Time `json:"recorded_at"`
    Issue       string    `json:"issue"`
    Processing  bool      `json:"processing"`
    ProcessedBy string    `json:"processed_by"`
    ProcessSteps []string `json:"process_steps"`
    Completed   bool      `json:"completed"`
    Result      string    `json:"result"`
}

var calls []Call

func main() {
    db, err := sql.Open("mysql", "root:mambau2001@tcp(127.0.0.1:3306)/calldb")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    r := gin.Default()
r.GET("/calls", getAllCalls)
r.POST("/calls", func(c *gin.Context) {
    var newCall Call
    if err := c.ShouldBindJSON(&newCall); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Check if RecordedAt is zero value and set to current time if so
    if newCall.RecordedAt.IsZero() {
        newCall.RecordedAt = time.Now()
    } else {
        recordedAtStr := newCall.RecordedAt.Format(time.RFC3339)
        recordedAtTime, err := time.Parse(time.RFC3339, recordedAtStr)
        if err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        newCall.RecordedAt = recordedAtTime
    }

    fmt.Println(newCall.RecordedAt)

    result, err := db.Exec("INSERT INTO calls (idcalls,parent_name, student_name, recorded_by, recorded_at, issue, processing, processed_by, process_steps, completed, result) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", 
        newCall.ParentName, newCall.StudentName, newCall.RecordedBy, newCall.RecordedAt, newCall.Issue, newCall.Processing, newCall.ProcessedBy, fmt.Sprintf("%v", newCall.ProcessSteps), newCall.Completed, newCall.Result)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    newCall.ID = strconv.FormatInt(id, 10)
    c.JSON(201, newCall)
})


    // Other endpoints...

    r.Run()
}

func recordCall(c *gin.Context) {
    var newCall Call
    if err := c.ShouldBindJSON(&newCall); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    newCall.ID = generateID()
    newCall.RecordedAt = time.Now()
    calls = append(calls, newCall)

    c.JSON(201, newCall)
}

func processCall(c *gin.Context) {
    id := c.Param("id")
    var processInfo struct {
        ProcessedBy  string   `json:"processed_by"`
        ProcessSteps []string `json:"process_steps"`
    }

    if err := c.ShouldBindJSON(&processInfo); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    for i, call := range calls {
        if call.ID == id {
            calls[i].Processing = true
            calls[i].ProcessedBy = processInfo.ProcessedBy
            calls[i].ProcessSteps = processInfo.ProcessSteps
            c.JSON(200, calls[i])
            return
        }
    }

    c.JSON(404, gin.H{"error": "Call not found"})
}

func completeCall(c *gin.Context) {
    id := c.Param("id")
    var completeInfo struct {
        Completed bool   `json:"completed"`
        Result    string `json:"result"`
    }

    if err := c.ShouldBindJSON(&completeInfo); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    for i, call := range calls {
        if call.ID == id {
            calls[i].Completed = completeInfo.Completed
            calls[i].Result = completeInfo.Result
            c.JSON(200, calls[i])
            return
        }
    }

    c.JSON(404, gin.H{"error": "Call not found"})
}

func getAllCalls(c *gin.Context) {
    db, err := sql.Open("mysql", "root:mambau2001@tcp(127.0.0.1:3306)/calldb")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer db.Close()

    rows, err := db.Query("SELECT * FROM calls")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    calls := []Call{}
    for rows.Next() {
        var call Call
        var recordedAtBytes []byte
        var processStepsBytes []byte
        err := rows.Scan(
            &call.ID, &call.ParentName, &call.StudentName, &call.RecordedBy, &recordedAtBytes,
            &call.Issue, &call.Processing, &call.ProcessedBy, &processStepsBytes, &call.Completed, &call.Result,
        )
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        if string(recordedAtBytes) == "0000-00-00" {
            call.RecordedAt = time.Time{}
        } else {
            recordedAtTime, err := time.Parse("2006-01-02 15:04:05", string(recordedAtBytes))
            if err != nil {
                c.JSON(500, gin.H{"error": err.Error()})
                return
            }
            call.RecordedAt = recordedAtTime
        }

        var processSteps []string
        err = json.Unmarshal(processStepsBytes, &processSteps)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        call.ProcessSteps = processSteps

        calls = append(calls, call)
    }

    if err := rows.Err(); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, calls)
}
func generateID() string {
    // Implement your own ID generation logic
    return "123456789"
}