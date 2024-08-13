package consumer

import (
	"context"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Npcf_UEPolicy"
	"github.com/free5gc/openapi/models"
)

// Request to PCF that create the UE policy
func UEPolicyControlCreate(ue *amf_context.AmfUe, anType models.AccessType) (*models.ProblemDetails, error) {
	logger.UePolicyLog.Info("Handle UEPolicyControlCreate !!")
	configuration := Npcf_UEPolicy.NewConfiguration()
	logger.UePolicyLog.Warnf("ue.PcfUri: %s", ue.PcfUri)
	configuration.SetBasePath(ue.PcfUri) //TODO:cheackbase path is "/npcf-ue-policy-control/v1/"
	client := Npcf_UEPolicy.NewAPIClient(configuration)

	amfSelf := amf_context.GetSelf()

	policyAssociationRequest := models.PolicyAssociationRequest{
		NotificationUri: amfSelf.GetIPv4Uri() + "/namf-callback/v1/ue-policy/",
		Supi:            ue.Supi,
		Pei:             ue.Pei,
		Gpsi:            ue.Gpsi,
		AccessType:      anType,
		ServingPlmn: &models.NetworkId{
			Mcc: ue.PlmnId.Mcc,
			Mnc: ue.PlmnId.Mnc,
		},
		GroupIds: ue.AccessAndMobilitySubscriptionData.InternalGroupIds,
		Guami:    &amfSelf.ServedGuamiList[0],
	}

	//TODO: Process Reponse and replace _ with res
	_, httpResp, localErr := client.DefaultApi.PoliciesPost(context.Background(), policyAssociationRequest)
	if localErr == nil {
		locationHeader := httpResp.Header.Get("Location")
		logger.ConsumerLog.Debugf("location header: %+v", locationHeader)
		ue.AmPolicyUri = locationHeader

		// re := regexp.MustCompile("/policies/.*")
		// match := re.FindStringSubmatch(locationHeader)

		// ue.PolicyAssociationId = match[0][10:]
		// ue.AmPolicyAssociation = &res

		// if res.Triggers != nil {
		// 	for _, trigger := range res.Triggers {
		// 		if trigger == models.RequestTrigger_LOC_CH {
		// 			ue.RequestTriggerLocationChange = true
		// 		}
		// 		//if trigger == models.RequestTrigger_PRA_CH {
		// 		// TODO: Presence Reporting Area handling (TS 23.503 6.1.2.5, TS 23.501 5.6.11)
		// 		//}
		// 	}
		// }

		// logger.ConsumerLog.Debugf("UE AM Policy Association ID: %s", ue.PolicyAssociationId)
		// logger.ConsumerLog.Debugf("AmPolicyAssociation: %+v", ue.AmPolicyAssociation)
		logger.UePolicyLog.Warnln("ue policy creation HTTP POST REQUEST no error!!")
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			return nil, localErr
		}
		problem := localErr.(Npcf_UEPolicy.GenericOpenAPIError).Model().(models.ProblemDetails)
		return &problem, nil
	} else {
		return nil, openapi.ReportError("server no response")
	}
	return nil, nil
}

func UEPolicyControlUpdate(ue *amf_context.AmfUe, updateRequest models.PolicyAssociationUpdateRequest) (
	problemDetails *models.ProblemDetails, err error) {
	configuration := Npcf_UEPolicy.NewConfiguration()
	configuration.SetBasePath(ue.PcfUri)
	client := Npcf_UEPolicy.NewAPIClient(configuration)

	//TODO: Process Reponse and replace _ with res
	_, httpResp, localErr := client.DefaultApi.PoliciesPolAssoIdUpdatePost(
		context.Background(), ue.PolicyAssociationId, updateRequest)
	if localErr == nil {
		// if res.ServAreaRes != nil {
		// 	ue.AmPolicyAssociation.ServAreaRes = res.ServAreaRes
		// }
		// if res.Rfsp != 0 {
		// 	ue.AmPolicyAssociation.Rfsp = res.Rfsp
		// }
		// ue.AmPolicyAssociation.Triggers = res.Triggers
		// ue.RequestTriggerLocationChange = false
		// for _, trigger := range res.Triggers {
		// 	if trigger == models.RequestTrigger_LOC_CH {
		// 		ue.RequestTriggerLocationChange = true
		// 	}
		// 	// if trigger == models.RequestTrigger_PRA_CH {
		// 	// TODO: Presence Reporting Area handling (TS 23.503 6.1.2.5, TS 23.501 5.6.11)
		// 	// }
		// }
		return
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("server no response")
	}
	return problemDetails, err
}

func UEPolicyControlDelete(ue *amf_context.AmfUe) (problemDetails *models.ProblemDetails, err error) {
	configuration := Npcf_UEPolicy.NewConfiguration()
	configuration.SetBasePath(ue.PcfUri)
	client := Npcf_UEPolicy.NewAPIClient(configuration)

	httpResp, localErr := client.DefaultApi.PoliciesPolAssoIdDelete(context.Background(), ue.PolicyAssociationId)
	if localErr == nil {
		ue.RemoveAmPolicyAssociation()
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("server no response")
	}

	return
}
