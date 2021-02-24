package infoblox

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func resourceCNAMERecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceCNAMERecordCreate,
		Read:   resourceCNAMERecordGet,
		Update: resourceCNAMERecordUpdate,
		Delete: resourceCNAMERecordDelete,

		Schema: map[string]*schema.Schema{
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone under which record has to be created.",
			},
			"dns_view": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Dns View under which the zone has been created.",
			},
			"canonical": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Canonical name for the record.",
			},
			"alias": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The alias name for the record.",
			},
			"vm_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance id.",
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of your tenant in cloud.",
			},
			"eas": &schema.Schema{
				Type: 	     schema.TypeMap,
				Optional:    true,
				Description: "Extensible Attributes"
			},
		},
	}
}
func resourceCNAMERecordCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning to create CNAME record ", resourceCNAMERecordIDString(d))

	zone := d.Get("zone").(string)
	dnsView := d.Get("dns_view").(string)
	canonical := d.Get("canonical").(string)
	alias := d.Get("alias").(string)
	if !strings.Contains(alias, zone) {
		alias = d.Get("alias").(string) + "." + zone
	}
	tenantID := d.Get("tenant_id").(string)
	vmId := d.Get("vm_id").(string)
	connector := m.(*ibclient.Connector)

	//eas := d.Get("eas").(map[string]interface{})
	
	eas := make(map[string]interface{})
	eas["activedirectory"] = d.Get("activedirectory")(string)
	eas["ITSMSR"] = d.Get("ITSMSR")(string)
	eas["country"] = d.Get("country")(string)
	eas["MPLSVRF"] = d.Get("MPLSVRF")(string)]
	eas["securityzone"] = d.Get("securityzone")(string)
	eas["cloudaccountemailDL"] = d.Get("cloudaccountemailDL")(string)
	eas["cloudaccountID"] = d.Get("cloudaccountID")(string)
	eas["cloudserviceprovider"] = d.Get("cloudserviceprovider")(string)
	eas["cloudzone"] = d.Get("cloudzone")(string)	
	eas["region"] = d.Get("region")(string)	
	eas["sitePOCemail"] = d.Get("sitePOCemail")(string)
	eas["passcode"] = d.Get("passcode")(string)
	
	ea := make(ibclient.EA)
	//for k,v := range eas {
		//ea[k] = v			
	//}	
	ea["activedirectory"] = eas["activedirectory"] 
	ea["ITSMSR"] = eas["ITSMSR"]
	ea["country"] = eas["country"]
	ea["MPLSVRF"] = eas["MPLSVRF"]
	ea["securityzone"] = eas["securityzone"]
	ea["cloudaccountemailDL"] = eas["cloudaccountemailDL"]
	ea["cloudaccountID"] = eas["cloudaccountID"]
	ea["cloudserviceprovider"] = eas["cloudserviceprovider"]
	ea["cloudzone"] = eas["cloudzone"]
	ea["region"] = eas["region"]
	ea["sitePOCemail"] = eas["sitePOCemail"]
	ea["passcode"] = eas["passcode"]
	
	ea["VM Name"] = canonical

	if vmId != "" {
		ea["VM ID"] = vmId
	}

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)
	recordCNAME, err := objMgr.CreateCNAMERecord(canonical, alias, dnsView, ea)
	if err != nil {
		return fmt.Errorf("Error creating CNAME Record : %s", err)
	}

	d.Set("recordName", alias)
	d.SetId(recordCNAME.Ref)

	log.Printf("[DEBUG] %s: Creation of CNAME Record complete", resourceCNAMERecordIDString(d))
	return resourceCNAMERecordGet(d, m)
}

func resourceCNAMERecordGet(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Begining to Get CNAME Record", resourceCNAMERecordIDString(d))

	dnsView := d.Get("dns_view").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	obj, err := objMgr.GetCNAMERecordByRef(d.Id())
	if err != nil {
		return fmt.Errorf("Getting CNAME RECORD failed from dns view(%s) : %s", dnsView, err)
	}
	d.SetId(obj.Ref)
	log.Printf("[DEBUG] %s: Completed reading required CNAME Record ", resourceCNAMERecordIDString(d))
	return nil
}

func resourceCNAMERecordUpdate(d *schema.ResourceData, m interface{}) error {

	return fmt.Errorf("updating CNAME record is not supported")
}

func resourceCNAMERecordDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Deletion of CNAME Record", resourceCNAMERecordIDString(d))

	dnsView := d.Get("dns_view").(string)
	tenantID := d.Get("tenant_id").(string)
	connector := m.(*ibclient.Connector)

	objMgr := ibclient.NewObjectManager(connector, "Terraform", tenantID)

	_, err := objMgr.DeleteCNAMERecord(d.Id())
	if err != nil {
		return fmt.Errorf("Deletion of CNAME Record failed with %s from dns view %s", dnsView, err)
	}
	d.SetId("")

	log.Printf("[DEBUG] %s: Deletion of CNAME Record complete", resourceCNAMERecordIDString(d))
	return nil
}

type resourceCNAMERecordIDStringInterface interface {
	Id() string
}

func resourceCNAMERecordIDString(d resourceCNAMERecordIDStringInterface) string {
	id := d.Id()
	if id == "" {
		id = "<new resource>"
	}
	return fmt.Sprintf("infoblox_cname_record (ID = %s)", id)
}
