package test

import (
	"context"

	"encoding/json"

	"time"

	"weavelab.xyz/deployer-resource-sync/pkg/resource"

	"weavelab.xyz/monorail/shared/wlib/werror"

	"weavelab.xyz/wstore/pkg/clusters"
)

type ClusterClient struct {
	*clusters.Client
}

func DefaultClusterClient(ctx context.Context) (*ClusterClient, error) {
	client, err := clusters.NewDefaultClient(ctx, clusters.ResourcesProjectID, clusters.ResourcesCollectionName)
	if err != nil {
		return nil, err
	}

	return &ClusterClient{client}, nil
}

func initClusterCache() error {
	clusterCache, shouldUpdate, err := checkClusterCache()
	if err != nil {
		return werror.Wrap(err, "failed to check if cluster cache needs update")
	}

	if !shouldUpdate {
		clusters.InitWithStatic(clusterCache)
		return nil
	}

	// update cache
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = clusters.Init(ctx)
	if err != nil {

		return werror.Wrap(err, "failed to initialize cluster cache")
	}

	return nil
}

func checkClusterCache() (clusterCache map[string]resource.Cluster, shouldUpdate bool, _ error) {
	if shouldUpdate {
		return nil, true, nil
	}

	err := json.Unmarshal(nil, &clusterCache)
	if err != nil {
		return clusterCache, true, werror.Wrap(err, "failed to unmarshal ~/.bart/clusters cache file")
	}

	if len(clusterCache) == 0 {
		return nil, true, nil
	}

	return clusterCache, shouldUpdate, nil
}
