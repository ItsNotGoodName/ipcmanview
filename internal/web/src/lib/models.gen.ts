/* Do not change, this code is generated from Golang structs */


export class WSData {
    type: string;
    data: any;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.type = source["type"];
        this.data = source["data"];
    }
}
export class WSEvent {
    action: string;
    data: Object;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.action = source["action"];
        this.data = source["data"];
    }
}
export class DahuaEvent {
    id: number;
    device_id: number;
    code: string;
    action: string;
    index: number;
    data: Object;
    created_at: Date;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.device_id = source["device_id"];
        this.code = source["code"];
        this.action = source["action"];
        this.index = source["index"];
        this.data = source["data"];
        this.created_at = new Date(source["created_at"]);
    }
}