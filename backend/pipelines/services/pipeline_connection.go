package services

import (
    "bytes"
    "strings"
    "net/http"
    "encoding/json"
    "ylem_pipelines/config"

    "github.com/kelseyhightower/envconfig"
)

func UpdatePipelineConnection(organizationUuid string, isPipelineCreated bool) bool {
    var config config.Config
    err := envconfig.Process("", &config)
    if err != nil {
        return false
    }

    url := strings.Replace(config.NetworkConfig.UpdateConnectionsUrl, "{uuid}", organizationUuid, -1);

    rp, _ := json.Marshal(map[string]bool{"is_pipeline_created": isPipelineCreated})
    var jsonStr = []byte(rp)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    if err != nil {
        return false
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        return true
    } else {
        return false
    }
}
