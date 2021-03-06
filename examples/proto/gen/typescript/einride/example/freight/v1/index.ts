export type Shipment = {
	name?: string;
	createTime?: unknown;
	updateTime?: unknown;
	deleteTime?: unknown;
	originSite?: string;
	destinationSite?: string;
	pickupEarliestTime?: unknown;
	pickupLatestTime?: unknown;
	deliveryEarliestTime?: unknown;
	deliveryLatestTime?: unknown;
	lineItems?: unknown[];
	annotations?: unknown[];
};

export type LineItem = {
	title?: string;
	quantity?: number;
	weightKg?: number;
	volumeM3?: number;
};

export type Shipper = {
	name?: string;
	createTime?: unknown;
	updateTime?: unknown;
	deleteTime?: unknown;
	displayName?: string;
};

export type Site = {
	name?: string;
	createTime?: unknown;
	updateTime?: unknown;
	deleteTime?: unknown;
	displayName?: string;
	latLng?: unknown;
};

export type GetShipperRequest = {
	name?: string;
};

export type ListShippersRequest = {
	pageSize?: number;
	pageToken?: string;
};

export type ListShippersResponse = {
	shippers?: unknown[];
	nextPageToken?: string;
};

export type CreateShipperRequest = {
	shipper?: unknown;
};

export type UpdateShipperRequest = {
	shipper?: unknown;
	updateMask?: unknown;
};

export type DeleteShipperRequest = {
	name?: string;
};

export type GetSiteRequest = {
	name?: string;
};

export type ListSitesRequest = {
	parent?: string;
	pageSize?: number;
	pageToken?: string;
};

export type ListSitesResponse = {
	sites?: unknown[];
	nextPageToken?: string;
};

export type CreateSiteRequest = {
	parent?: string;
	site?: unknown;
};

export type UpdateSiteRequest = {
	site?: unknown;
	updateMask?: unknown;
};

export type DeleteSiteRequest = {
	name?: string;
};

export type GetShipmentRequest = {
	name?: string;
};

export type ListShipmentsRequest = {
	parent?: string;
	pageSize?: number;
	pageToken?: string;
};

export type ListShipmentsResponse = {
	shipments?: unknown[];
	nextPageToken?: string;
};

export type CreateShipmentRequest = {
	parent?: string;
	shipment?: unknown;
};

export type UpdateShipmentRequest = {
	shipment?: unknown;
	updateMask?: unknown;
};

export type DeleteShipmentRequest = {
	name?: string;
};

