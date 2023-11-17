create table tb_application
(
    id                         int auto_increment
        primary key,
    application_name           varchar(64)                         not null comment '应用名称',
    application_administrators int                                 null comment 'app管理员',
    application_type           varchar(8)                          not null comment 'app类型 WEB | APPLICATION',
    application_path           varchar(128)                        not null comment '应用路径（默认应用名称）',
    must_contain_language      json                                null comment '必须包含的语言',
    application_environment    varchar(16)                         not null comment '系统环境 STG & DEV & PROD',
    dual_authentication        int       default 0                 null comment '是否开启双重认证 0关闭，1开启',
    create_time                timestamp default CURRENT_TIMESTAMP null comment '创建时间',
    create_user_id             int                                 null comment '创建人用户ID',
    update_time                timestamp default CURRENT_TIMESTAMP null comment '更新时间',
    update_user_id             int                                 null comment '更新用户ID',
    constraint tb_application_application_name_uindex
        unique (application_name),
    constraint tb_application_pk
        unique (application_path)
)
    comment '注册应用表' row_format = DYNAMIC;

create table tb_application_globalization_document_code
(
    document_id              int auto_increment
        primary key,
    application_id           int                                 not null comment '应用ID',
    namespace_id             int                                 not null,
    document_code            varchar(255)                        not null comment '文案编码',
    document_desc            varchar(128)                        null comment '文案描述',
    is_enable                int       default 1                 not null comment '是否上线',
    online_time              timestamp                           null comment '上线时间',
    online_operator_user_id  int                                 null comment '上线操作人',
    offline_time             timestamp                           null comment '下线时间',
    offline_operator_user_id int                                 null comment '下线操作人',
    offline_access_user_id   int                                 null comment '下线审核人',
    create_time              timestamp default CURRENT_TIMESTAMP null comment '创建时间',
    update_time              timestamp                           null on update CURRENT_TIMESTAMP comment '更新时间',
    create_user_id           int                                 null comment '创建人',
    delete_flag              int       default 0                 null comment '删除标识',
    delete_time              timestamp                           null comment '删除时间',
    delete_user_id           int                                 null comment '删除操作人',
    remarks                  varchar(255)                        null comment '备注',
    constraint tb_application_globalization_biz_uindex
        unique (application_id, namespace_id, document_code),
    constraint tb_application_globalization_document_code_uindex
        unique (namespace_id, application_id, document_code)
)
    comment '应用多语言' row_format = DYNAMIC;

create table tb_application_globalization_document_value
(
    id                   int auto_increment comment 'PK'
        primary key,
    document_id          int                                 null comment '文案编码ID',
    namespace_id         int                                 not null comment '命名空间ID',
    country_iso          varchar(2)                          null comment '国家二字码',
    country_name         varchar(32)                         null comment '国家名称',
    document_value       varchar(5120)                       null comment '文言',
    document_is_online   int       default 1                 not null comment '文案上线状态',
    create_time          timestamp default CURRENT_TIMESTAMP null comment '文案创建时间',
    create_user_id       int                                 null comment '创建人',
    update_time          timestamp                           null on update CURRENT_TIMESTAMP comment '更新时间',
    update_user_id       int                                 null comment '更新人',
    last_update_document varchar(5120)                       null comment '上一次更新文案',
    delete_flag          int                                 null comment '删除时间',
    delete_time          timestamp                           null comment '删除时间',
    delete_user_id       int                                 null comment '删除人'
)
    comment '文案值' row_format = DYNAMIC;

create index country_iso_index
    on tb_application_globalization_document_value (country_iso);

create index document_id_index
    on tb_application_globalization_document_value (document_id);

create index namespace_id_index
    on tb_application_globalization_document_value (namespace_id);

create table tb_application_namespace
(
    namespace_id        int auto_increment
        primary key,
    namespace_code      varchar(64)                         null comment '命名空间编码',
    namespace_name      varchar(64)                         not null comment '命名空间名称',
    namespace_path      varchar(256)                        not null comment '命名空间路径',
    namespace_parent_id int                                 null comment '父命名空间ID',
    application_id      int                                 not null comment '关联的应用ID',
    create_time         timestamp default CURRENT_TIMESTAMP null comment '创建时间',
    create_user         int                                 null comment '创建人',
    constraint tb_application_namespace_namespace_code_uindex
        unique (namespace_code),
    constraint tb_application_namespace_namespace_id_namespace_path_uindex
        unique (namespace_id, namespace_path),
    constraint tb_application_namespace_namespace_path_application_id_uindex
        unique (namespace_path, application_id)
)
    comment '应用命名空间' row_format = DYNAMIC;

