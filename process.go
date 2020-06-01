package cmon

func (client *Client) GetSqlProcesses(clusterID uint64) ([]*SQLProcess, error) {
	req := &GetSqlProcessesRequest{
		WithOperation: &WithOperation{"getSqlProcesses"},
		WithClusterID: &WithClusterID{clusterID},
	}
	res := &GetSqlProcessesResponse{}
	if err := client.Request(ModuleClusters, req, res, false); err != nil {
		return nil, err
	}
	if res.RequestStatus != RequestStatusOk {
		return nil, NewErrorFromResponseData(res.WithResponseData)
	}
	return res.Processes, nil
}

type GetSqlProcessesRequest struct {
	*WithOperation `json:",inline"`
	*WithClusterID `json:",inline"`
}

type GetSqlProcessesResponse struct {
	*WithResponseData `json:",inline"`

	Processes []*SQLProcess `json:"processes"`
}

type SQLProcess struct {
	BlockedByTrxID string `json:"blocked_by_trx_id"`
	Client         string `json:"client"`
	Command        string `json:"command"`
	CurrentTime    int64  `json:"currentTime"`
	DB             string `json:"db"`
	Duration       int64  `json:"duration"`
	Host           string `json:"host"`
	HostID         uint64 `json:"host_id"`
	Hostname       string `json:"hostname"`
	Info           string `json:"info"`
	InnodbStatus   string `json:"innodb_status"`
	InnodbTrxID    string `json:"innodb_trx_id"`
	Instance       string `json:"instance"`
	Lastseen       int64  `json:"lastseen"`
	Message        string `json:"message"`
	MysqlTrxID     int64  `json:"mysql_trx_id"`
	PID            int64  `json:"pid"`
	Query          string `json:"query"`
	ReportTS       int64  `json:"report_ts"`
	SQL            string `json:"sql"`
	State          string `json:"state"`
	Time           int64  `json:"time"`
	User           string `json:"user"`
}
