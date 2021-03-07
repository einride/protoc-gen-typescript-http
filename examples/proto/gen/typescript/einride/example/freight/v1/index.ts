// A shipment represents transportation of goods between an origin
// [site][einride.example.freight.v1.Site] and a destination
// [site][einride.example.freight.v1.Site].
export type Shipment = {
	// The resource name of the shipment.
	name?: string;
	// The creation timestamp of the shipment.
	createTime?: wellKnownTimestamp;
	// The last update timestamp of the shipment.
	// Updated when create/update/delete operation is shipment.
	updateTime?: wellKnownTimestamp;
	// The deletion timestamp of the shipment.
	deleteTime?: wellKnownTimestamp;
	// The resource name of the origin site of the shipment.
	// Format: shippers/{shipper}/sites/{site}
	originSite?: string;
	// The resource name of the destination site of the shipment.
	// Format: shippers/{shipper}/sites/{site}
	destinationSite?: string;
	// The earliest pickup time of the shipment at the origin site.
	pickupEarliestTime?: wellKnownTimestamp;
	// The latest pickup time of the shipment at the origin site.
	pickupLatestTime?: wellKnownTimestamp;
	// The earliest delivery time of the shipment at the destination site.
	deliveryEarliestTime?: wellKnownTimestamp;
	// The latest delivery time of the shipment at the destination site.
	deliveryLatestTime?: wellKnownTimestamp;
	// The line items of the shipment.
	lineItems?: LineItem[];
	// Annotations of the shipment.
	annotations?: { [key: string]: string};
};

// Encoded using RFC 3339, where generated output will always be Z-normalized
// and uses 0, 3, 6 or 9 fractional digits.
// Offsets other than "Z" are also accepted.
type wellKnownTimestamp = string;

// A shipment line item.
export type LineItem = {
	// The title of the line item.
	title?: string;
	// The quantity of the line item.
	quantity?: number;
	// The weight of the line item in kilograms.
	weightKg?: number;
	// The volume of the line item in cubic meters.
	volumeM3?: number;
};

// A shipper is a supplier or owner of goods to be transported.
export type Shipper = {
	// The resource name of the shipper.
	name?: string;
	// The creation timestamp of the shipper.
	createTime?: wellKnownTimestamp;
	// The last update timestamp of the shipper.
	// Updated when create/update/delete operation is performed.
	updateTime?: wellKnownTimestamp;
	// The deletion timestamp of the shipper.
	deleteTime?: wellKnownTimestamp;
	// The display name of the shipper.
	displayName?: string;
};

// A site is a node in a [shipper][einride.example.freight.v1.Shipper]'s
// transport network.
export type Site = {
	// The resource name of the site.
	name?: string;
	// The creation timestamp of the site.
	createTime?: wellKnownTimestamp;
	// The last update timestamp of the site.
	// Updated when create/update/delete operation is performed.
	updateTime?: wellKnownTimestamp;
	// The deletion timestamp of the site.
	deleteTime?: wellKnownTimestamp;
	// The display name of the site.
	displayName?: string;
	// The geographic location of the site.
	latLng?: googletype_LatLng;
};

// An object that represents a latitude/longitude pair. This is expressed as a
// pair of doubles to represent degrees latitude and degrees longitude. Unless
// specified otherwise, this must conform to the
// <a href="http://www.unoosa.org/pdf/icg/2012/template/WGS_84.pdf">WGS84
// standard</a>. Values must be within normalized ranges.
export type googletype_LatLng = {
	// The latitude in degrees. It must be in the range [-90.0, +90.0].
	latitude?: number;
	// The longitude in degrees. It must be in the range [-180.0, +180.0].
	longitude?: number;
};

// Request message for FreightService.GetShipper.
export type GetShipperRequest = {
	// The resource name of the shipper to retrieve.
	// Format: shippers/{shipper}
	name?: string;
};

// Request message for FreightService.ListShippers.
export type ListShippersRequest = {
	// Requested page size. Server may return fewer shippers than requested.
	// If unspecified, server will pick an appropriate default.
	pageSize?: number;
	// A token identifying a page of results the server should return.
	// Typically, this is the value of
	// [ListShippersResponse.next_page_token][einride.example.freight.v1.ListShippersResponse.next_page_token]
	// returned from the previous call to `ListShippers` method.
	pageToken?: string;
};

// Response message for FreightService.ListShippers.
export type ListShippersResponse = {
	// The list of shippers.
	shippers?: Shipper[];
	// A token to retrieve next page of results.  Pass this value in the
	// [ListShippersRequest.page_token][einride.example.freight.v1.ListShippersRequest.page_token]
	// field in the subsequent call to `ListShippers` method to retrieve the next
	// page of results.
	nextPageToken?: string;
};

// Request message for FreightService.CreateShipper.
export type CreateShipperRequest = {
	// The shipper to create.
	shipper?: Shipper;
};

// Request message for FreightService.UpdateShipper.
export type UpdateShipperRequest = {
	// The shipper to update with. The name must match or be empty.
	// The shipper's `name` field is used to identify the shipper to be updated.
	// Format: shippers/{shipper}
	shipper?: Shipper;
	// The list of fields to be updated.
	updateMask?: wellKnownFieldMask;
};

// In JSON, a field mask is encoded as a single string where paths are
// separated by a comma. Fields name in each path are converted
// to/from lower-camel naming conventions.
// As an example, consider the following message declarations:
//
//     message Profile {
//       User user = 1;
//       Photo photo = 2;
//     }
//     message User {
//       string display_name = 1;
//       string address = 2;
//     }
//
// In proto a field mask for `Profile` may look as such:
//
//     mask {
//       paths: "user.display_name"
//       paths: "photo"
//     }
//
// In JSON, the same mask is represented as below:
//
//     {
//       mask: "user.displayName,photo"
//     }
type wellKnownFieldMask = string;

// Request message for FreightService.DeleteShipper.
export type DeleteShipperRequest = {
	// The resource name of the shipper to delete.
	// Format: shippers/{shipper}
	name?: string;
};

// Request message for FreightService.GetSite.
export type GetSiteRequest = {
	// The resource name of the site to retrieve.
	// Format: shippers/{shipper}/sites/{site}
	name?: string;
};

// Request message for FreightService.ListSites.
export type ListSitesRequest = {
	// The resource name of the parent, which owns this collection of sites.
	// Format: shippers/{shipper}
	parent?: string;
	// Requested page size. Server may return fewer sites than requested.
	// If unspecified, server will pick an appropriate default.
	pageSize?: number;
	// A token identifying a page of results the server should return.
	// Typically, this is the value of
	// [ListSitesResponse.next_page_token][einride.example.freight.v1.ListSitesResponse.next_page_token]
	// returned from the previous call to `ListSites` method.
	pageToken?: string;
};

// Response message for FreightService.ListSites.
export type ListSitesResponse = {
	// The list of sites.
	sites?: Site[];
	// A token to retrieve next page of results.  Pass this value in the
	// [ListSitesRequest.page_token][einride.example.freight.v1.ListSitesRequest.page_token]
	// field in the subsequent call to `ListSites` method to retrieve the next
	// page of results.
	nextPageToken?: string;
};

// Request message for FreightService.CreateSite.
export type CreateSiteRequest = {
	// The resource name of the parent shipper for which this site will be created.
	// Format: shippers/{shipper}
	parent?: string;
	// The site to create.
	site?: Site;
};

// Request message for FreightService.UpdateSite.
export type UpdateSiteRequest = {
	// The site to update with. The name must match or be empty.
	// The site's `name` field is used to identify the site to be updated.
	// Format: shippers/{shipper}/sites/{site}
	site?: Site;
	// The list of fields to be updated.
	updateMask?: wellKnownFieldMask;
};

// Request message for FreightService.DeleteSite.
export type DeleteSiteRequest = {
	// The resource name of the site to delete.
	// Format: shippers/{shipper}/sites/{site}
	name?: string;
};

// Request message for FreightService.GetShipment.
export type GetShipmentRequest = {
	// The resource name of the shipment to retrieve.
	// Format: shippers/{shipper}/shipments/{shipment}
	name?: string;
};

// Request message for FreightService.ListShipments.
export type ListShipmentsRequest = {
	// The resource name of the parent, which owns this collection of shipments.
	// Format: shippers/{shipper}
	parent?: string;
	// Requested page size. Server may return fewer shipments than requested.
	// If unspecified, server will pick an appropriate default.
	pageSize?: number;
	// A token identifying a page of results the server should return.
	// Typically, this is the value of
	// [ListShipmentsResponse.next_page_token][einride.example.freight.v1.ListShipmentsResponse.next_page_token]
	// returned from the previous call to `ListShipments` method.
	pageToken?: string;
};

// Response message for FreightService.ListShipments.
export type ListShipmentsResponse = {
	// The list of shipments.
	shipments?: Shipment[];
	// A token to retrieve next page of results.  Pass this value in the
	// [ListShipmentsRequest.page_token][einride.example.freight.v1.ListShipmentsRequest.page_token]
	// field in the subsequent call to `ListShipments` method to retrieve the next
	// page of results.
	nextPageToken?: string;
};

// Request message for FreightService.CreateShipment.
export type CreateShipmentRequest = {
	// The resource name of the parent shipper for which this shipment will be created.
	// Format: shippers/{shipper}
	parent?: string;
	// The shipment to create.
	shipment?: Shipment;
};

// Request message for FreightService.UpdateShipment.
export type UpdateShipmentRequest = {
	// The shipment to update with. The name must match or be empty.
	// The shipment's `name` field is used to identify the shipment to be updated.
	// Format: shippers/{shipper}/shipments/{shipment}
	shipment?: Shipment;
	// The list of fields to be updated.
	updateMask?: wellKnownFieldMask;
};

// Request message for FreightService.DeleteShipment.
export type DeleteShipmentRequest = {
	// The resource name of the shipment to delete.
	// Format: shippers/{shipper}/shipments/{shipment}
	name?: string;
};

export interface FreightService {
	GetShipper(request: GetShipperRequest): Promise<Shipper>
	ListShippers(request: ListShippersRequest): Promise<ListShippersResponse>
	CreateShipper(request: CreateShipperRequest): Promise<Shipper>
	UpdateShipper(request: UpdateShipperRequest): Promise<Shipper>
	DeleteShipper(request: DeleteShipperRequest): Promise<Shipper>
	GetSite(request: GetSiteRequest): Promise<Site>
	ListSites(request: ListSitesRequest): Promise<ListSitesResponse>
	CreateSite(request: CreateSiteRequest): Promise<Site>
	UpdateSite(request: UpdateSiteRequest): Promise<Site>
	DeleteSite(request: DeleteSiteRequest): Promise<Site>
	GetShipment(request: GetShipmentRequest): Promise<Shipment>
	ListShipments(request: ListShipmentsRequest): Promise<ListShipmentsResponse>
	CreateShipment(request: CreateShipmentRequest): Promise<Shipment>
	UpdateShipment(request: UpdateShipmentRequest): Promise<Shipment>
	DeleteShipment(request: DeleteShipmentRequest): Promise<Shipment>
}

type requestHandler = (path: string, method: string, body: string | null) => Promise<unknown>

export function createFreightServiceClient(handler: requestHandler): FreightService {
	return {
		GetShipper(request) {
			return handler("", "", null) as Promise<Shipper>
		},
		ListShippers(request) {
			return handler("", "", null) as Promise<ListShippersResponse>
		},
		CreateShipper(request) {
			return handler("", "", null) as Promise<Shipper>
		},
		UpdateShipper(request) {
			return handler("", "", null) as Promise<Shipper>
		},
		DeleteShipper(request) {
			return handler("", "", null) as Promise<Shipper>
		},
		GetSite(request) {
			return handler("", "", null) as Promise<Site>
		},
		ListSites(request) {
			return handler("", "", null) as Promise<ListSitesResponse>
		},
		CreateSite(request) {
			return handler("", "", null) as Promise<Site>
		},
		UpdateSite(request) {
			return handler("", "", null) as Promise<Site>
		},
		DeleteSite(request) {
			return handler("", "", null) as Promise<Site>
		},
		GetShipment(request) {
			return handler("", "", null) as Promise<Shipment>
		},
		ListShipments(request) {
			return handler("", "", null) as Promise<ListShipmentsResponse>
		},
		CreateShipment(request) {
			return handler("", "", null) as Promise<Shipment>
		},
		UpdateShipment(request) {
			return handler("", "", null) as Promise<Shipment>
		},
		DeleteShipment(request) {
			return handler("", "", null) as Promise<Shipment>
		},
	}
}
