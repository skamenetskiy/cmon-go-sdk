package cmon

func (client *Client) GetAllClusterInfo(req *GetAllClusterInfoRequest) ([]*Cluster, error) {
	if req == nil {
		req = &GetAllClusterInfoRequest{
			WithOperation: &WithOperation{"getAllClusterInfo"},
		}
	} else {
		req.WithOperation = &WithOperation{"getAllClusterInfo"}
	}
	res := &GetAllClusterInfoResponse{}
	if err := client.Request(ModuleClusters, req, res, false); err != nil {
		return nil, err
	}
	if res.RequestStatus != RequestStatusOk {
		return nil, NewErrorFromResponseData(res.WithResponseData)
	}
	return res.Clusters, nil
}

type GetAllClusterInfoRequest struct {
	*WithOperation `json:",inline"`

	WithHosts        bool `json:"with_hosts"`
	WithSheetInfo    bool `json:"with_sheet_info"`
	WithDatabases    bool `json:"with_databases"`
	WithLicenseCheck bool `json:"with_license_check"`
	WithTags         bool `json:"with_tags"`
}

type GetAllClusterInfoResponse struct {
	*WithResponseData `json:",inline"`
	*WithTotal        `json:",inline"`

	Clusters []*Cluster
}

type Cluster struct {
	*WithClassName `json:",inline"`

	ClusterID   uint64 `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	ClusterType string `json:"cluster_type"`
}
