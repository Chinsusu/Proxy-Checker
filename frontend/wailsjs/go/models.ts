export namespace models {
	
	export class IPQualityResult {
	    ip: string;
	    port?: string;
	    status: string;
	    country: string;
	    city: string;
	    region: string;
	    vpn: boolean;
	    proxy: boolean;
	    isp: string;
	    organization: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new IPQualityResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.status = source["status"];
	        this.country = source["country"];
	        this.city = source["city"];
	        this.region = source["region"];
	        this.vpn = source["vpn"];
	        this.proxy = source["proxy"];
	        this.isp = source["isp"];
	        this.organization = source["organization"];
	        this.error = source["error"];
	    }
	}
	export class WhoisResult {
	    ip: string;
	    country: string;
	    country_code: string;
	    region: string;
	    city: string;
	    flag: string;
	    isp: string;
	    asn: string;
	    timezone: string;
	    status: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new WhoisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.country = source["country"];
	        this.country_code = source["country_code"];
	        this.region = source["region"];
	        this.city = source["city"];
	        this.flag = source["flag"];
	        this.isp = source["isp"];
	        this.asn = source["asn"];
	        this.timezone = source["timezone"];
	        this.status = source["status"];
	        this.error = source["error"];
	    }
	}

}

