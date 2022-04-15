package resolver

//	protocolVersion 协议号
const V_1_0 string = "1.0"
const V_1_1 string = "1.1"
const V_2_0 string = "2.0"
const V_2_1 string = "2.1"
const V_2_2 string = "2.2"

const KEY_DATA string = "data"

const KEY_SIGN string = "sign"

const KEY_TOPIC string = "topic"

const KEY_PAYLOAD string = "payload"

const KEY_PROTOCOL string = "protocol"

const KEY_DATA_TIME string = "t"

const KEY_DATA_POINT string = "dps"

const SLASH string = "/"

const KEY_USERNAME string = "username"

const MSG_ID string = "msg_id"

const KEY_TIME_STAMP string = "ts"

const KEY_EVENT_TYPE string = "etype"

const KEY_EVENT_DATA string = "edata"

const KEY_SUCCESS string = "success"

const KEY_DEVID string = "devId"

const KEY_GWID string = "gwId"

const KEY_LOG_URL string = "logUrl"

const KEY_UUID string = "uuid"

const KEY_BUCKET string = "bucket"

const KEY_LOG_SUF string = "logSuf"

const KEY_SN_DATA string = "sn"

const KEY_CT string = "ct"

const KEY_S string = "s"

const KEY_TYPE string = "type"

const GATEWAY_ADD_DEL string = "0"                    //网关增删，APP订阅"),
const GATEWAY_UPDATE string = "1"                     //网关更新，APP订阅"),
const DEVICE_ADD_DEL string = "2"                     //设备增删，APP订阅"),
const DEVICE_UPDATE string = "3"                      //设备更新，APP订阅"),
const DATA_POINT_REPORT string = "4"                  //数据点上报，APP订阅"),
const HARDWARE_UPGRADE_MOBILE string = "9"            //固件升级状态接收，APP订阅"),
const DATA_POINT_ISSUE string = "5"                   //数据点下发，GW订阅"),
const DEVICE_DEL_PUSH_GATEWAY string = "8"            //设备删除通知原网关，GW订阅"),
const HARDWARE_UPGRADE string = "10"                  //固件升级确认升级，GW订阅"),
const GATEWAY_DEL_PUSH_GATEWAY string = "11"          //app发起网关重置操作通知网关重置自身，GW订阅"),
const DEVICE_TIMER_CHANGED string = "13"              //网关定时钟发生变化通知网关，GW订阅"),
const APP_DEVICE_TIMER_CHANGED string = "14"          //网关定时钟发生变化通知网关，APP订阅"),
const HARDWARE_UPGRADE_V2 string = "15"               //固件升级确认升级新版本，GW订阅"),
const HARDWARE_UPGRADE_PROGRESS string = "16"         //固件升级进度，APP订阅"),
const DEVICE_REPORT_SELF_INFO string = "17"           //设备上报状态信息，云端订阅"),
const DEVICE_REQ_CLOUD_DATA string = "18"             //设备向云端请求数据，云端订阅"),
const DEVICE_RES_CLOUD_DATA string = "19"             //设备向云端请求数据云端返回响应，GW订阅"),
const DEVICE_UPDATE_TO_3RD_CLOUD string = "20"        //设备信息发生变化推送给第三方云，第三方云订阅"),
const DEVICE_SWITCH_REGION string = "21"              //设备切换可用区，GW订阅"),
const SEND_REQUEST_TO_DEVICE string = "22"            //向设备发起查询请求，目前主要是APP在用，设备订阅"),
const RECEIVE_RESPONSE_FROM_DEVICE string = "23"      //设备回应查询请求，目前主要APP在使用，APP订阅"),
const REQUEST_DEVICE_DP_REPORT string = "24"          //向设备发起查询设备dp点当前状态，设备订阅"),
const SUB_DEVICE_ONLINE_STATUS_REPORT string = "25"   //上报子设备上下线状态，云端订阅"),
const DEVICE_RES_ACK string = "26"                    //下发指令，设备收到指令后返回ACK，云端订阅"),
const DEVICE_DYNAMIC_CONFIG_UPDATE string = "27"      //设备配置信息发生变化，需要从云端重新获取配置数据，设备订阅"),
const DEVICE_LEVEL_DP_RAW_PUBLISH string = "28"       //APP给设备发送透传指令，设备端订阅"),
const DEVICE_LEVEL_DP_RAW_REPORT string = "29"        //设备上报透传指令，APP订阅"),
const BATCH_DATA_POINT_REPORT string = "30"           //多个设备数据点上报，APP订阅"),
const DEVICE_ONLINE_OFFLINE string = "32"             //连接mqtt设备的在线离线状态，APP订阅, 设备topic"),
const SUB_DEVICE_ADD_DEL string = "33"                //子设备增删除，APP订阅，设备订阅, 设备topic"),
const DEVICE_NAME_UPDATE string = "34"                //设备改名，APP订阅, 设备topic"),
const MESH_ADD_DEL string = "35"                      //MESH增删，APP订阅, 组topic"),
const DEVICE_GROUP_UPDATE string = "37"               //设备群组新增/编辑/解散"),
const LOCATION_ADD_DEL string = "39"                  //家庭增删，APP订阅, 用户topic"),
const LOCATION_INFO_UPDATE string = "40"              //家庭信息更新，APP订阅, 组topic"),
const DST_INTERVALS_UPDATE string = "41"              //夏令时区间更新，GW订阅"),
const DEVICE_EVENT_NOTIFY string = "43"               //设备发生事件通知APP和云端"),
const DEVICE_VERSION_UPDATE string = "44"             //设备版本更新，APP订阅，设备topic"),
const DEVICE_LOG_UPLOAD string = "45"                 //设备日志上传，云端订阅，设备topic"),
const GW_SUB_BIND_ENABLE string = "200"               //app使能网关添加子设备请求"),
const D_GROUP_DP_PUB string = "47"                    //群组dp下发"),
const D_GROUP_DP_NAME_UPDATE string = "53"            //群组dp改名"),
const SUB_DEVICE_TRANSFER string = "54"               //子设备转移，app订阅"),
const DEVICE_PUBLISH_EVENT string = "55"              //云端下发给网关事件 ，设备订阅"),
const DEVICE_BIND_INFO_NOTIFY string = "46"           //云端下发硬件配网通知消息"),
const SHARING_EVENT_OPEN_NOTIFY string = "48"         //分享事件对外通知"),
const GATEWAY_NOTIFY string = "49"                    //向网关发送通知"),
const DP_NAME_UPDATE string = "50"                    //dp名称修改推送给app"),
const CLOUD_EVENT_NOTIFY string = "52"                //云端发事件通知APP"),
const NOTIFY_DEVICE_BIND_INFO_BY_DEVICE string = "51" //已配网设备向未配网设备发送ssid,pwd,token，通知配网"),
const CLOUD_ALARM_TO_APP string = "56"                //云端产品告警通知APP"),
const REPORT_DEVICE_INFO string = "57"                //上报查询设备信息"),
const PUBLIC_DEVICE_INFO string = "58"                //下发查询设备信息"),
const LOCAL_SCENE_EXECUTE string = "206"              //执行zigbee本地场景"),
const CAMERA_ORDER_EVENT string = "300"               //摄像头订单事件,服务端下发,设备订阅"),
const IPC_RTSP_NOTIFY string = "301"                  //ipc rtsp服务器等待请求通知"),
const IPC_P2P_SIGNAL string = "302"                   //ipc p2p呼叫信号"),
const FACE_DETECT_UPDATE string = "303"               //摄像头唤醒指令"),
const FACE_DETECT_SAMPLE_UPDATE string = "304"        //人脸检测样本数据更新 string = "新增/删除/变更)"),
const SENCE_EVNET string = "400"                      //场景相关事件"),
const SENCE_CURD string = "401"                       //场景增/删/查/改"),
const LOG string = "999"                              //日志"),
const SUB_DEVICE_GROUP_LOCAL_CURD string = "202"      //子设备本地组操作 添加/删除"),
const SUB_DEVICE_GROUP_CLOUD_CURD string = "203"      //子设备组操作通知 云端通知app"),
const FACE_DOOR_SYNC string = "306"                   //人脸识别门禁数据同步通知"),
const VOICE_EVENT string = "501"                      //语音相关事件"),
const IPC_EXPAND_CMD string = "307"                   //摄像头 云端扩展指令 string = "非dp指令)下发"),
const SECURITY_SYS_CORRELATION_PUBLISH string = "701" //安防相关下发"),
const SECURITY_SYS_CORRELATION_REPORT string = "702"  //安防相关上报");
