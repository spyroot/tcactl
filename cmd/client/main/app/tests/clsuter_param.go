package tests

//
//var spec request.ClusterMgmt
//spec.Name = "test123"
//spec.ClusterType = string(request.ClusterManagement)
//spec.ClusterTemplateId = "46e2f7c5-1908-4bee-9a58-8808ff57a2e2"
//spec.HcxCloudUrl = "https://tca-cp03.cnfdemo.io"
//spec.VmTemplate = "/Datacenter/vm/templates/photon-3-kube-v1.20.4+vmware.1"
//spec.EndpointIP = "10.247.7.212"
//spec.PlacementParams = []models.PlacementParams{
//		*models.NewPlacementParams("templates", "Folder"),
//		*models.NewPlacementParams("vsanDatastore", "Datastore"),
//		*models.NewPlacementParams("pod03", "ResourcePool"),
//		*models.NewPlacementParams("core", "IsValidClusterCompute"),
//}
//spec.ClusterPassword = "Vk13YXJlMSE"
//
//net := models.NewNetworks(string(request.ClusterManagement),
//	"/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
//	[]string{"10.246.2.9"})
//
//master := models.NewTypeNode("master",  []models.Networks{*net}, []models.PlacementParams{
//	*models.NewPlacementParams("Discovered virtual machine", "Folder"),
//	*models.NewPlacementParams("vsanDatastore", "Datastore"),
//	*models.NewPlacementParams("pod03", "ResourcePool"),
//	*models.NewPlacementParams("core", "IsValidClusterCompute"),
//})
//
//worker := models.NewTypeNode("note-pool01", []models.Networks{*net}, []models.PlacementParams{
//	*models.NewPlacementParams("Discovered virtual machine", "Folder"),
//	*models.NewPlacementParams("vsanDatastore", "Datastore"),
//	*models.NewPlacementParams("pod03", "ResourcePool"),
//	*models.NewPlacementParams("core", "IsValidClusterCompute"),
//})
//
//
//spec.MasterNodes = []models.TypeNode{*master}
//spec.WorkerNodes =  []models.TypeNode{*worker}
//
//
//var spec request.ClusterMgmt
