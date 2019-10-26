package apicupcfg

type ManagementSubsysDescriptor interface {
	GetManagementSubsysName() string
	GetPlatformApiEndpoint() string
	GetConsumerApiEndpoint() string
	GetApiManagerUIEndpoint() string
	GetCloudAdminUIEndpoint() string
}

func (mgt *MgtSubsysVm) GetManagementSubsysName() string { return mgt.SubsysName }
func (mgt *MgtSubsysVm) GetPlatformApiEndpoint() string { return mgt.PlatformApi }
func (mgt *MgtSubsysVm) GetConsumerApiEndpoint() string { return mgt.ConsumerApi }
func (mgt *MgtSubsysVm) GetApiManagerUIEndpoint() string { return mgt.ApiManagerUi }
func (mgt *MgtSubsysVm) GetCloudAdminUIEndpoint() string { return mgt.CloudAdminUi }

func (mgt *MgtSubsysK8s) GetManagementSubsysName() string { return mgt.SubsysName }
func (mgt *MgtSubsysK8s) GetPlatformApiEndpoint() string { return mgt.PlatformApi }
func (mgt *MgtSubsysK8s) GetConsumerApiEndpoint() string { return mgt.ConsumerApi }
func (mgt *MgtSubsysK8s) GetApiManagerUIEndpoint() string { return mgt.ApiManagerUi }
func (mgt *MgtSubsysK8s) GetCloudAdminUIEndpoint() string { return mgt.CloudAdminUi }

type PortalSubsysDescriptor interface {
	GetPortalSubsysName() string
	GetPortalWWWEndpoint() string
	GetPortalAdminEndpoint() string
}

func (ptl *PtlSubsysVm) GetPortalSubsysName() string { return ptl.SubsysName }
func (ptl *PtlSubsysVm) GetPortalWWWEndpoint() string { return ptl.PortalWww }
func (ptl *PtlSubsysVm) GetPortalAdminEndpoint() string { return ptl.PortalAdmin }

func (ptl *PtlSubsysK8s) GetPortalSubsysName() string { return ptl.SubsysName }
func (ptl *PtlSubsysK8s) GetPortalWWWEndpoint() string { return ptl.PortalWWW }
func (ptl *PtlSubsysK8s) GetPortalAdminEndpoint() string { return ptl.PortalAdmin }

type AnalyticsSubsysDescriptor interface {
	GetAnalyticsSubsysName() string
	GetAnalyticsIngestionEndpoint() string
	GetAnalyticsClientEndpoint() string
}

func (alt *AltSubsysVm) GetAnalyticsSubsysName() string { return alt.SubsysName }
func (alt *AltSubsysVm) GetAnalyticsIngestionEndpoint() string { return alt.AnalyticsIngestion }
func (alt *AltSubsysVm) GetAnalyticsClientEndpoint() string { return alt.AnalyticsClient }

func (alt *AlytSubsysK8s) GetAnalyticsSubsysName() string { return alt.SubsysName }
func (alt *AlytSubsysK8s) GetAnalyticsIngestionEndpoint() string { return alt.AnalyticsIngestionEndpoint }
func (alt *AlytSubsysK8s) GetAnalyticsClientEndpoint() string { return alt.AnalyticsClientEndpoint }

type GatewaySubsysDescriptor interface {
	GetGatewaySubsysName() string
	GetApicGatewayServiceEndpoint() string
	GetApiGatewayEndpoint() string
}

func (gwy *GwySubsysVm) GetGatewaySubsysName() string { return gwy.SubsysName }
func (gwy *GwySubsysVm) GetApicGatewayServiceEndpoint() string { return gwy.ApicGwService }
func (gwy *GwySubsysVm) GetApiGatewayEndpoint() string { return gwy.ApiGateway }

func (gwy *GwSubsysK8s) GetGatewaySubsysName() string { return gwy.SubsysName }
func (gwy *GwSubsysK8s) GetApicGatewayServiceEndpoint() string { return gwy.ApicGwService }
func (gwy *GwSubsysK8s) GetApiGatewayEndpoint() string { return gwy.ApiGateway }
