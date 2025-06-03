import 'package:get/get.dart';
import 'package:json_annotation/json_annotation.dart';

part 'weibo.g.dart';

@JsonSerializable(genericArgumentFactories: true)
class WeiboResponse<T> {
  WeiboResponse({this.data, required this.ok, this.msg});

  T? data;
  int ok;
  String? msg;

  factory WeiboResponse.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic json) fromJsonT,
  ) => _$WeiboResponseFromJson(json, fromJsonT);

  Map<String, dynamic> toJson(T Function(T) toJsonT) =>
      _$WeiboResponseToJson(this, toJsonT);
}

@JsonSerializable()
class WeiboList {
  @JsonKey(name: 'bottom_tips_text')
  String bottomTipsText;
  @JsonKey(name: 'bottom_tips_visible')
  bool bottomTipsVisible;
  List<Myblog> list;
  @JsonKey(name: 'since_id')
  String sinceId;
  @JsonKey(name: 'status_visible')
  int statusVisible;
  List<dynamic> topicList;
  int total;

  WeiboList({
    required this.bottomTipsText,
    required this.bottomTipsVisible,
    required this.list,
    required this.sinceId,
    required this.statusVisible,
    required this.topicList,
    required this.total,
  });

  factory WeiboList.fromJson(Map<String, dynamic> json) =>
      _$WeiboListFromJson(json);

  Map<String, dynamic> toJson() => _$WeiboListToJson(this);
}


@JsonSerializable()
class WeiboOriginList {
  List<Myblog> list;
  dynamic total;

  WeiboOriginList({
    required this.list,
    required this.total,
  });

  factory WeiboOriginList.fromJson(Map<String, dynamic> json) =>
      _$WeiboOriginListFromJson(json);

  Map<String, dynamic> toJson() => _$WeiboOriginListToJson(this);
}


@JsonSerializable()
class Myblog {
  @JsonKey(name: 'analysis_extra')
  String? analysisExtra;
  List<dynamic>? annotations;
  @JsonKey(name: 'attitudes_count')
  int? attitudesCount;
  @JsonKey(name: 'attitudes_status')
  int attitudesStatus;
  @JsonKey(name: 'can_edit')
  bool? canEdit;
  @JsonKey(name: 'comment_manage_info')
  CommentManageInfo? commentManageInfo;
  @JsonKey(name: 'comments_count')
  int? commentsCount;
  @JsonKey(name: 'content_auth')
  int? contentAuth;
  @JsonKey(name: 'created_at')
  String createdAt;
  bool? favorited;
  int id;
  String idstr;
  @JsonKey(name: 'is_long_text')
  bool? isLongText;
  @JsonKey(name: 'is_paid')
  bool? isPaid;
  @JsonKey(name: 'is_show_bulletin')
  int? isShowBulletin;
  @JsonKey(name: 'is_show_mixed')
  bool? isShowMixed;
  bool isSinglePayAudio;
  @JsonKey(name: 'mblog_vip_type')
  int? mblogVipType;
  String mblogid;
  int? mblogtype;
  String mid;
  @JsonKey(name: 'mixed_count')
  int? mixedCount;
  int? mlevel;
  @JsonKey(name: 'number_display_strategy')
  NumberDisplayStrategy? numberDisplayStrategy;
  @JsonKey(name: 'page_info')
  PageInfo? pageInfo;
  @JsonKey(name: 'pic_ids')
  List<String>? picIds;
  @JsonKey(name: 'pic_infos')
  Map<String, AllPicInfo2>? picInfos;
  @JsonKey(name: 'pic_num')
  int? picNum;
  bool pictureViewerSign;
  List<dynamic> rcList;
  String? readtimetype;
  @JsonKey(name: 'region_name')
  String? regionName;
  @JsonKey(name: 'repost_type')
  int? repostType;
  @JsonKey(name: 'reposts_count')
  int? repostsCount;
  @JsonKey(name: 'retweeted_status')
  Myblog? retweetedStatus;
  String? rid;
  @JsonKey(name: 'share_repost_type')
  int? shareRepostType;
  bool showFeedComment;
  @JsonKey(name: 'show_feed_repost')
  bool? showFeedRepost;
  bool showPictureViewer;
  String source;
  String text;
  int textLength;
  @JsonKey(name: 'text_raw')
  String textRaw;
  @JsonKey(name: 'url_struct')
  List<UrlStruct>? urlStruct;
  User? user;
  Visible visible;

  Myblog({
    this.analysisExtra,
    this.annotations,
    this.attitudesCount,
    required this.attitudesStatus,
    this.canEdit,
    this.commentManageInfo,
    this.commentsCount,
    this.contentAuth,
    required this.createdAt,
    this.favorited,
    required this.id,
    required this.idstr,
    this.isLongText,
    this.isPaid,
    this.isShowBulletin,
    this.isShowMixed,
    required this.isSinglePayAudio,
    this.mblogVipType,
    required this.mblogid,
    this.mblogtype,
    required this.mid,
    this.mixedCount,
    this.mlevel,
    this.numberDisplayStrategy,
    this.pageInfo,
    this.picIds,
    this.picInfos,
    this.picNum,
    required this.pictureViewerSign,
    required this.rcList,
    this.readtimetype,
    this.regionName,
    this.repostType,
    this.repostsCount,
    this.retweetedStatus,
    this.rid,
    this.shareRepostType,
    required this.showFeedComment,
    this.showFeedRepost,
    required this.showPictureViewer,
    required this.source,
    required this.text,
    required this.textLength,
    required this.textRaw,
    this.urlStruct,
    this.user,
    required this.visible,
  });

  factory Myblog.fromJson(Map<String, dynamic> json) => _$MyblogFromJson(json);

  Map<String, dynamic> toJson() => _$MyblogToJson(this);
}

@JsonSerializable()
class PageInfo {
  @JsonKey(name: 'act_status')
  int actStatus;
  Actionlog actionlog;
  @JsonKey(name: 'author_id')
  String authorId;

  String authorid;
  String content1;
  String content2;
  @JsonKey(name: 'media_info')
  MediaInfo mediaInfo;
  @JsonKey(name: 'object_id')
  String objectId;
  @JsonKey(name: 'object_type')
  String objectType;
  String oid;
  @JsonKey(name: 'page_id')
  String pageId;
  @JsonKey(name: 'page_pic')
  String pagePic;
  @JsonKey(name: 'page_title')
  String pageTitle;
  @JsonKey(name: 'page_url')
  String pageUrl;
  @JsonKey(name: 'pic_info')
  AllPicInfo picInfo;
  @JsonKey(name: 'short_url')
  String shortUrl;
  String type;
  @JsonKey(name: 'type_icon')
  String typeIcon;
  String warn;

  PageInfo({
    required this.actStatus,
    required this.actionlog,
    required this.authorId,
    required this.authorid,
    required this.content1,
    required this.content2,
    required this.mediaInfo,
    required this.objectId,
    required this.objectType,
    required this.oid,
    required this.pageId,
    required this.pagePic,
    required this.pageTitle,
    required this.pageUrl,
    required this.picInfo,
    required this.shortUrl,
    required this.type,
    required this.typeIcon,
    required this.warn,
  });

  factory PageInfo.fromJson(Map<String, dynamic> json) =>
      _$PageInfoFromJson(json);

  Map<String, dynamic> toJson() => _$PageInfoToJson(this);
}

@JsonSerializable()
class AllPicInfo {
  @JsonKey(name: 'pic_big')
  PicInfo picBig;
  @JsonKey(name: 'pic_middle')
  PicInfo picMiddle;
  @JsonKey(name: 'pic_small')
  PicInfo picSmall;

  AllPicInfo({
    required this.picBig,
    required this.picMiddle,
    required this.picSmall,
  });

  factory AllPicInfo.fromJson(Map<String, dynamic> json) =>
      _$AllPicInfoFromJson(json);

  Map<String, dynamic> toJson() => _$AllPicInfoToJson(this);
}

@JsonSerializable()
class PicInfo {
  String? height;
  String? url;
  String? width;

  PicInfo({this.height, this.url, this.width});

  factory PicInfo.fromJson(Map<String, dynamic> json) =>
      _$PicInfoFromJson(json);

  Map<String, dynamic> toJson() => _$PicInfoToJson(this);
}

@JsonSerializable()
class Actionlog {
  @JsonKey(name: 'act_code')
  dynamic actCode;
  @JsonKey(name: 'act_type')
  int? actType;
  String? ext;
  String? fid;
  String? lcardid;
  String? mid;
  String oid;
  String? source;
  int? uuid;

  Actionlog({
    required this.actCode,
     this.actType,
    this.ext,
    this.fid,
    this.lcardid,
    this.mid,
    required this.oid,
    this.source,
    this.uuid,
  });

  factory Actionlog.fromJson(Map<String, dynamic> json) =>
      _$ActionlogFromJson(json);

  Map<String, dynamic> toJson() => _$ActionlogToJson(this);
}

@JsonSerializable()
class MediaInfo {
  @JsonKey(name: 'act_status')
  int actStatus;
  @JsonKey(name: 'author_info')
  User authorInfo;
  @JsonKey(name: 'author_mid')
  String authorMid;
  @JsonKey(name: 'author_name')
  String authorName;
  @JsonKey(name: 'belong_collection')
  int belongCollection;
  @JsonKey(name: 'big_pic_info')
  PicInfo bigPicInfo;
  int duration;
  @JsonKey(name: 'ext_info')
  ExtInfo extInfo;
  @JsonKey(name: 'extra_info')
  ExtraInfo extraInfo;
  String format;
  @JsonKey(name: 'forward_strategy')
  int forwardStrategy;
  @JsonKey(name: 'h265_mp4_hd')
  String? h265Mp4Hd;
  @JsonKey(name: 'h265_mp4_ld')
  String? h265Mp4Ld;
  @JsonKey(name: 'h5_url')
  String? h5Url;
  @JsonKey(name: 'hevc_mp4_hd')
  String? hevcMp4720p;
  @JsonKey(name: 'hevc_mp4_ld')
  String? inch4Mp4Hd;
  @JsonKey(name: 'inch_5_5_mp4_hd')
  String? inch55Mp4Hd;
  @JsonKey(name: 'inch_5_mp4_hd')
  String? inch5Mp4Hd;
  @JsonKey(name: 'is_keep_current_mblog')
  int? isKeepCurrentMblog;
  @JsonKey(name: 'is_short_video')
  int isShortVideo;
  @JsonKey(name: 'jump_to')
  int jumpTo;
  @JsonKey(name: 'kol_title')
  String kolTitle;
  @JsonKey(name: 'media_id')
  String mediaId;
  @JsonKey(name: 'mp4_720p_mp4')
  String mp4720pMp4;
  @JsonKey(name: 'mp4_hd_url')
  String mp4HdUrl;
  @JsonKey(name: 'mp4_sd_url')
  String mp4SdUrl;
  @JsonKey(name: 'name')
  String name;
  @JsonKey(name: 'next_title')
  String nextTitle;
  @JsonKey(name: 'online_users')
  String onlineUsers;
  @JsonKey(name: 'online_users_number')
  int onlineUsersNumber;
  @JsonKey(name: 'origin_total_bitrate')
  int originTotalBitrate;
  @JsonKey(name: 'play_completion_actions')
  List<PlayCompletionAction> playCompletionActions;
  @JsonKey(name: 'play_loop_type')
  int playLoopType;
  @JsonKey(name: 'playback_list')
  List<Playback>? playbackList;
  @JsonKey(name: 'prefetch_size')
  int prefetchSize;
  @JsonKey(name: 'prefetch_type')
  int prefetchType;
  String protocol;
  @JsonKey(name: 'search_scheme')
  String searchScheme;
  @JsonKey(name: 'show_mute_button')
  bool showMuteButton;
  @JsonKey(name: 'show_progress_bar')
  int showProgressBar;
  @JsonKey(name: 'storage_type')
  String storageType;
  @JsonKey(name: 'stream_url')
  String streamUrl;
  @JsonKey(name: 'stream_url_hd')
  String streamUrlHd;
  @JsonKey(name: 'titles_display_time')
  String titlesDisplayTime;
  int ttl;
  @JsonKey(name: 'video_download_strategy')
  VideoDownloadStrategy videoDownloadStrategy;
  @JsonKey(name: 'video_orientation')
  String videoOrientation;
  @JsonKey(name: 'video_publish_time')
  int videoPublishTime;
  @JsonKey(name: 'vote_is_show')
  int voteIsShow;

  MediaInfo({
    required this.actStatus,
    required this.authorInfo,
    required this.authorMid,
    required this.authorName,
    required this.belongCollection,
    required this.bigPicInfo,
    required this.duration,
    required this.extInfo,
    required this.extraInfo,
    required this.format,
    required this.forwardStrategy,
    this.h265Mp4Hd,
    this.h265Mp4Ld,
    this.h5Url,
    this.hevcMp4720p,
    this.inch4Mp4Hd,
    this.inch55Mp4Hd,
    this.inch5Mp4Hd,
    this.isKeepCurrentMblog,
    required this.isShortVideo,
    required this.jumpTo,
    required this.kolTitle,
    required this.mediaId,
    required this.mp4720pMp4,
    required this.mp4HdUrl,
    required this.mp4SdUrl,
    required this.name,
    required this.nextTitle,
    required this.onlineUsers,
    required this.onlineUsersNumber,
    required this.originTotalBitrate,
    required this.playCompletionActions,
    required this.playLoopType,
    this.playbackList,
    required this.prefetchSize,
    required this.prefetchType,
    required this.protocol,
    required this.searchScheme,
    required this.showMuteButton,
    required this.showProgressBar,
    required this.storageType,
    required this.streamUrl,
    required this.streamUrlHd,
    required this.titlesDisplayTime,
    required this.ttl,
    required this.videoDownloadStrategy,
    required this.videoOrientation,
    required this.videoPublishTime,
    required this.voteIsShow,
  });

  factory MediaInfo.fromJson(Map<String, dynamic> json) =>
      _$MediaInfoFromJson(json);

  Map<String, dynamic> toJson() => _$MediaInfoToJson(this);
}

@JsonSerializable()
class StatusTotalCounter {
  @JsonKey(name: 'comment_cnt')
  String commentCnt;
  @JsonKey(name: 'like_cnt')
  String likeCnt;
  @JsonKey(name: 'repost_cnt')
  String repostCnt;
  @JsonKey(name: 'total_cnt')
  String totalCnt;
  @JsonKey(name: 'total_cnt_format')
  dynamic totalCntFormat;

  StatusTotalCounter({
    required this.commentCnt,
    required this.likeCnt,
    required this.repostCnt,
    required this.totalCnt,
    required this.totalCntFormat,
  });

  factory StatusTotalCounter.fromJson(Map<String, dynamic> json) =>
      _$StatusTotalCounterFromJson(json);

  Map<String, dynamic> toJson() => _$StatusTotalCounterToJson(this);
}

@JsonSerializable()
class PlayCompletionAction {
  Actionlog actionlog;
  @JsonKey(name: 'btn_code')
  int? btnCode;
  @JsonKey(name: 'countdown_time')
  int? countdownTime;
  @JsonKey(name: 'display_endtime')
  int? displayEndtime;
  @JsonKey(name: 'display_mode')
  int? displayMode;
  @JsonKey(name: 'display_start_time')
  int? displayStartTime;
  @JsonKey(name: 'display_type')
  int? displayType;
  Ext? ext;
  String? icon;
  String? link;
  String? scheme;
  @JsonKey(name: 'show_position')
  int? showPosition;
  String? text;
  dynamic type;

  PlayCompletionAction({
    required this.actionlog,
     this.btnCode,
    this.countdownTime,
    this.displayEndtime,
    this.displayMode,
    this.displayStartTime,
    this.displayType,
    this.ext,
     this.icon,
     this.link,
    this.scheme,
     this.showPosition,
     this.text,
     this.type,
  });

  factory PlayCompletionAction.fromJson(Map<String, dynamic> json) =>
      _$PlayCompletionActionFromJson(json);

  Map<String, dynamic> toJson() => _$PlayCompletionActionToJson(this);
}

@JsonSerializable()
class Ext {
  @JsonKey(name: 'followers_count')
  int? followersCount;
  int? level;
  String? uid;
  @JsonKey(name: 'user_name')
  String? userName;
  bool? verified;
  @JsonKey(name: 'verified_reason')
  String? verifiedReason;
  @JsonKey(name: 'verified_type')
  int? verifiedType;

  Ext({
    this.followersCount,
     this.level,
     this.uid,
     this.userName,
    this.verified,
    this.verifiedReason,
     this.verifiedType,
  });

  factory Ext.fromJson(Map<String, dynamic> json) => _$ExtFromJson(json);

  Map<String, dynamic> toJson() => _$ExtToJson(this);
}

class ExtraInfo {
  String sceneid;

  ExtraInfo({required this.sceneid});

  factory ExtraInfo.fromJson(Map<String, dynamic> json) {
    return ExtraInfo(sceneid: json['sceneid']);
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['sceneid'] = sceneid;
    return data;
  }
}

class Playback {
  Meta meta;
  PlayInfo playInfo;

  Playback({required this.meta, required this.playInfo});

  factory Playback.fromJson(Map<String, dynamic> json) {
    return Playback(
      meta: Meta.fromJson(json['meta']),
      playInfo: PlayInfo.fromJson(json['play_info']),
    );
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['meta'] = meta.toJson();
    data['play_info'] = playInfo.toJson();
    return data;
  }
}

@JsonSerializable()
class PlayInfo {
  @JsonKey(name: 'audio_bits_per_sample')
  int? audioBitsPerSample;
  @JsonKey(name: 'audio_channels')
  int? audioChannels;
  @JsonKey(name: 'audio_codecs')
  String? audioCodecs;
  @JsonKey(name: 'audio_sample_fmt')
  String? audioSampleFmt;
  @JsonKey(name: 'audio_sample_rate')
  int? audioSampleRate;
  int? bitrate;
  @JsonKey(name: 'color_transfer')
  String? colorTransfer;
  @JsonKey(name: 'dolby_atmos')
  bool? dolbyAtmos;
  double? duration;
  Extension extension;
  @JsonKey(name: 'first_pkt_end_pos')
  int? firstPktEndPos;
  int? fps;
  int height;
  String label;
  String mime;
  @JsonKey(name: 'prefetch_enabled')
  bool prefetchEnabled;
  @JsonKey(name: 'prefetch_range')
  String prefetchRange;
  String protocol;
  @JsonKey(name: 'quality_class')
  String qualityClass;
  @JsonKey(name: 'quality_desc')
  String qualityDesc;
  @JsonKey(name: 'quality_label')
  String qualityLabel;
  String? sar;
  int? size;
  @JsonKey(name: 'stereo_video')
  int? stereoVideo;
  @JsonKey(name: 'tcp_receive_buffer')
  int tcpReceiveBuffer;
  int type;
  String url;
  @JsonKey(name: 'video_codecs')
  String? videoCodecs;
  @JsonKey(name: 'video_decoder')
  String videoDecoder;
  String? watermark;
  int width;

  PlayInfo({
    this.audioBitsPerSample,
    this.audioChannels,
    this.audioCodecs,
    this.audioSampleFmt,
    this.audioSampleRate,
    this.bitrate,
    this.colorTransfer,
    this.dolbyAtmos,
    this.duration,
    required this.extension,
    this.firstPktEndPos,
    this.fps,
    required this.height,
    required this.label,
    required this.mime,
    required this.prefetchEnabled,
    required this.prefetchRange,
    required this.protocol,
    required this.qualityClass,
    required this.qualityDesc,
    required this.qualityLabel,
    this.sar,
    this.size,
    this.stereoVideo,
    required this.tcpReceiveBuffer,
    required this.type,
    required this.url,
    this.videoCodecs,
    required this.videoDecoder,
    this.watermark,
    required this.width,
  });

  factory PlayInfo.fromJson(Map<String, dynamic> json) =>
      _$PlayInfoFromJson(json);

  Map<String, dynamic> toJson() => _$PlayInfoToJson(this);
}

class Extension {
  TranscodeInfo transcodeInfo;

  Extension({required this.transcodeInfo});

  factory Extension.fromJson(Map<String, dynamic> json) {
    return Extension(
      transcodeInfo: TranscodeInfo.fromJson(json['transcode_info']),
    );
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['transcode_info'] = transcodeInfo.toJson();
    return data;
  }
}

@JsonSerializable()
class TranscodeInfo {
  @JsonKey(name: 'ab_strategies')
  String abStrategies;
  @JsonKey(name: 'origin_video_dr')
  String originVideoDr;
  @JsonKey(name: 'pcdn_jank')
  int pcdnJank;
  @JsonKey(name: 'pcdn_rule_id')
  int pcdnRuleId;

  TranscodeInfo({
    required this.abStrategies,
    required this.originVideoDr,
    required this.pcdnJank,
    required this.pcdnRuleId,
  });

  factory TranscodeInfo.fromJson(Map<String, dynamic> json) =>
      _$TranscodeInfoFromJson(json);

  Map<String, dynamic> toJson() => _$TranscodeInfoToJson(this);
}

@JsonSerializable()
class Meta {
  @JsonKey(name: 'is_hidden')
  bool isHidden;
  String label;
  @JsonKey(name: 'quality_class')
  String qualityClass;
  @JsonKey(name: 'quality_desc')
  String qualityDesc;
  @JsonKey(name: 'quality_group')
  int qualityGroup;
  @JsonKey(name: 'quality_index')
  int qualityIndex;
  @JsonKey(name: 'quality_label')
  String qualityLabel;
  int type;

  Meta({
    required this.isHidden,
    required this.label,
    required this.qualityClass,
    required this.qualityDesc,
    required this.qualityGroup,
    required this.qualityIndex,
    required this.qualityLabel,
    required this.type,
  });

  factory Meta.fromJson(Map<String, dynamic> json) => _$MetaFromJson(json);

  Map<String, dynamic> toJson() => _$MetaToJson(this);
}

class ExtInfo {
  String videoOrientation;

  ExtInfo({required this.videoOrientation});

  factory ExtInfo.fromJson(Map<String, dynamic> json) {
    return ExtInfo(videoOrientation: json['video_orientation']);
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['video_orientation'] = videoOrientation;
    return data;
  }
}

class VideoDownloadStrategy {
  int abandonDownload;

  VideoDownloadStrategy({required this.abandonDownload});

  factory VideoDownloadStrategy.fromJson(Map<String, dynamic> json) {
    return VideoDownloadStrategy(abandonDownload: json['abandon_download']);
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['abandon_download'] = abandonDownload;
    return data;
  }
}

@JsonSerializable()
class CommentManageInfo {
  @JsonKey(name: 'approval_comment_type')
  int approvalCommentType;
  @JsonKey(name: 'comment_permission_type')
  int commentPermissionType;
  @JsonKey(name: 'comment_sort_type')
  int commentSortType;
  @JsonKey(name: 'ai_play_picture_type')
  int? aiPlayPictureType;

  CommentManageInfo({
    required this.approvalCommentType,
    required this.commentPermissionType,
    required this.commentSortType,
    this.aiPlayPictureType,
  });

  factory CommentManageInfo.fromJson(Map<String, dynamic> json) =>
      _$CommentManageInfoFromJson(json);

  Map<String, dynamic> toJson() => _$CommentManageInfoToJson(this);
}

@JsonSerializable()
class UrlStruct {
  Actionlog actionlog;
  @JsonKey(name: 'h5_target_url')
  String h5TargetUrl;
  int hide;
  @JsonKey(name: 'long_url')
  String longUrl;
  @JsonKey(name: 'need_save_obj')
  int needSaveObj;
  @JsonKey(name: 'object_type')
  String objectType;
  @JsonKey(name: 'ori_url')
  String oriUrl;
  @JsonKey(name: 'page_id')
  String pageId;
  bool result;
  @JsonKey(name: 'short_url')
  String shortUrl;
  @JsonKey(name: 'storage_type')
  String storageType;
  int? ttl;
  @JsonKey(name: 'url_title')
  String urlTitle;
  @JsonKey(name: 'url_type')
  int urlType;
  @JsonKey(name: 'url_type_pic')
  String urlTypePic;

  UrlStruct({
    required this.actionlog,
    required this.h5TargetUrl,
    required this.hide,
    required this.longUrl,
    required this.needSaveObj,
    required this.objectType,
    required this.oriUrl,
    required this.pageId,
    required this.result,
    required this.shortUrl,
    required this.storageType,
    this.ttl,
    required this.urlTitle,
    required this.urlType,
    required this.urlTypePic,
  });

  factory UrlStruct.fromJson(Map<String, dynamic> json) =>
      _$UrlStructFromJson(json);

  Map<String, dynamic> toJson() => _$UrlStructToJson(this);
}

@JsonSerializable()
class User {
  @JsonKey(name: 'avatar_hd')
  String avatarHd;
  @JsonKey(name: 'avatar_large')
  String avatarLarge;
  String domain;
  @JsonKey(name: 'follow_me')
  bool followMe;
  bool following;
  @JsonKey(name: 'icon_list')
  List<dynamic>? iconList;
  int id;
  String idstr;
  int mbrank;
  int mbtype;
  @JsonKey(name: 'pc_new')
  int pcNew;
  @JsonKey(name: 'planet_video')
  bool planetVideo;
  @JsonKey(name: 'profile_image_url')
  String profileImageUrl;
  @JsonKey(name: 'profile_url')
  String profileUrl;
  @JsonKey(name: 'screen_name')
  String screenName;
  @JsonKey(name: 'status_total_counter')
  StatusTotalCounter? statusTotalCounter;
  @JsonKey(name: 'user_ability')
  int userAbility;
  @JsonKey(name: 'v_plus')
  int vPlus;
  bool verified;
  @JsonKey(name: 'verified_type')
  int verifiedType;
  String weihao;
  @JsonKey(name: 'verified_type_ext')
  int? verifiedTypeExt;
  @JsonKey(name: 'cover_image_phone')
  String? coverImagePhone;
  String? description;
  @JsonKey(name: 'followers_count')
  int? followersCount;
  @JsonKey(name: 'followers_count_str')
  String? followersCountStr;
  @JsonKey(name: 'friends_count')
  int? friendsCount;
  String? gender;

  String? location;
  @JsonKey(name: 'statuses_count')
  int? statusesCount;
  int? svip;
  String? url;
  @JsonKey(name: 'verified_reason')
  String? verifiedReason;

  int? vvip;

  User({
    required this.avatarHd,
    required this.avatarLarge,
    required this.domain,
    required this.followMe,
    required this.following,
    this.iconList,
    required this.id,
    required this.idstr,
    required this.mbrank,
    required this.mbtype,
    required this.pcNew,
    required this.planetVideo,
    required this.profileImageUrl,
    required this.profileUrl,
    required this.screenName,
     this.statusTotalCounter,
    required this.userAbility,
    required this.vPlus,
    required this.verified,
    required this.verifiedType,
    required this.weihao,
    this.verifiedTypeExt,
    this.coverImagePhone,
    this.description,
    this.followersCount,
    this.followersCountStr,
    this.friendsCount,
    this.gender,
    this.location,
    this.statusesCount,
    this.svip,
    this.url,
    this.verifiedReason,
    this.vvip,
  });

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);

  Map<String, dynamic> toJson() => _$UserToJson(this);
}

class NumberDisplayStrategy {
  int applyScenarioFlag;
  String displayText;
  int displayTextMinNumber;

  NumberDisplayStrategy({
    required this.applyScenarioFlag,
    required this.displayText,
    required this.displayTextMinNumber,
  });

  factory NumberDisplayStrategy.fromJson(Map<String, dynamic> json) {
    return NumberDisplayStrategy(
      applyScenarioFlag: json['apply_scenario_flag'],
      displayText: json['display_text'],
      displayTextMinNumber: json['display_text_min_number'],
    );
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['apply_scenario_flag'] = applyScenarioFlag;
    data['display_text'] = displayText;
    data['display_text_min_number'] = displayTextMinNumber;
    return data;
  }
}

class Visible {
  int listId;
  int type;

  Visible({required this.listId, required this.type});

  factory Visible.fromJson(Map<String, dynamic> json) {
    return Visible(listId: json['list_id'], type: json['type']);
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['list_id'] = listId;
    data['type'] = type;
    return data;
  }
}

@JsonSerializable()
class AllPicInfo2 {
  PicInfo2 bmiddle;
  String? fid;
  PicInfo2 large;
  PicInfo2 largecover;
  PicInfo2 largest;
  PicInfo2 mw2000;
  @JsonKey(name: 'object_id')
  String objectId;
  PicInfo2 original;
  @JsonKey(name: 'photo_tag')
  int photoTag;
  @JsonKey(name: 'pic_id')
  String picId;
  @JsonKey(name: 'pic_status')
  int picStatus;
  PicInfo2 thumbnail;
  String type;
  String? video;

  AllPicInfo2({
    required this.bmiddle,
    this.fid,
    required this.large,
    required this.largecover,
    required this.largest,
    required this.mw2000,
    required this.objectId,
    required this.original,
    required this.photoTag,
    required this.picId,
    required this.picStatus,
    required this.thumbnail,
    required this.type,
    this.video,
  });

  factory AllPicInfo2.fromJson(Map<String, dynamic> json) =>
      _$AllPicInfo2FromJson(json);

  Map<String, dynamic> toJson() => _$AllPicInfo2ToJson(this);
}

@JsonSerializable()
class PicInfo2 {
  @JsonKey(name: 'cut_type')
  int cutType;
  int height;
  String? type;
  String url;
  int width;

  PicInfo2({
    required this.cutType,
    required this.height,
     this.type,
    required this.url,
    required this.width,
  });

  factory PicInfo2.fromJson(Map<String, dynamic> json) =>
      _$PicInfo2FromJson(json);

  Map<String, dynamic> toJson() => _$PicInfo2ToJson(this);
}

@JsonSerializable()
class Icon {
  IconData data;
  String type;

  Icon({required this.data, required this.type});

  factory Icon.fromJson(Map<String, dynamic> json) => _$IconFromJson(json);

  Map<String, dynamic> toJson() => _$IconToJson(this);
}

@JsonSerializable()
class IconData {
  int mbrank;
  int mbtype;
  int svip;
  int vvip;

  IconData({
    required this.mbrank,
    required this.mbtype,
    required this.svip,
    required this.vvip,
  });

  factory IconData.fromJson(Map<String, dynamic> json) =>
      _$IconDataFromJson(json);

  Map<String, dynamic> toJson() => _$IconDataToJson(this);
}

@JsonSerializable()
class Button {
  Actionlog actionlog;
  String name;
  Params params;
  String type;

  Button({
    required this.actionlog,
    required this.name,
    required this.params,
    required this.type,
  });

  factory Button.fromJson(Map<String, dynamic> json) => _$ButtonFromJson(json);

  Map<String, dynamic> toJson() => _$ButtonToJson(this);
}

@JsonSerializable()
class Params {
  @JsonKey(name: 'disable_group')
  int disableGroup;
  Extparams extparams;
  int uid;

  Params({
    required this.disableGroup,
    required this.extparams,
    required this.uid,
  });

  factory Params.fromJson(Map<String, dynamic> json) => _$ParamsFromJson(json);

  Map<String, dynamic> toJson() => _$ParamsToJson(this);
}

@JsonSerializable()
class Extparams {
  String followcardid;

  Extparams({required this.followcardid});

  factory Extparams.fromJson(Map<String, dynamic> json) =>
      _$ExtparamsFromJson(json);

  Map<String, dynamic> toJson() => _$ExtparamsToJson(this);
}
