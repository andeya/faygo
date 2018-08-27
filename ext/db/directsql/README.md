# directSQL 使用说明

## 升级
   2018.08.26
      <parameter>的 default 增加2个默认参数
             - int64uuid ：64位整数长度的唯一id （通过配置机器节点支持分布式唯一id生成）
             - shortuuid ：短字符串的唯一id（将int64uuid转为36进制的值（10个数字+26个字母组成的））
     2018.08.22
     sql节点配置增加 eachtran 属性，配置 batchexec、batchmultiexec类型生成的所有SQL是否一个事务中执行，默认为false，true的则每个批次循环在不同的事务。
     batchexec 类型SQL支持配置多个sql    
     2017.08
    增加两种类型的sql，处理二进制对象保存到数据库和从数据库获取
    - getblob: ST_GETBLOB 获取BLOB (binary large object)，二进制大对象从数据库
    - setblob: ST_SETBLOB 保存BLOB (binary large object)，二进制大对象到数据库
   
## 简介
    directSQL通过配置sql映射文件，配置路由后可直接执行配置的sql并返回结果到客户端。

## 直接执行SQL引擎
    - 通过url直接执行SQL，参数采用POST方式传递  /bos/... sql标示
    - sqlengine--API执行sql的引擎
    - sqlhandle--处理执行请求的handle(API)
    - sqlrouter--SQL路由处理
    - sqlmanage--管理维护SQL定义文件到内存中
    - sqlservice--通过代码调用的接口单元
    - checkparameter -- 参数处理（参数验证与默认参数处理）
    - resultcache --查询结果缓存处理单元
    - sqlcontext --sql参数默认值的自定义函数单元
    - sqlhelper--辅助函数
    - sqlwatcher--SQL配置文件监控自动更新(  实现文件修改，删除监控，改名，新增貌似不行)
    - 系统中通过代码如何调用：
        directsql/sqlService 单元中的函数

## 配置文件：
    config/directsql.ini 的内容
     `
    ;SQL配置文件扩展名，只能一个。
    ext=.msql 
    ;是否懒惰加载，true 第一次请求才加载请求的sqlmodel文件然后缓存(未实现)，false＝一开始就根据配置的roots目录全部加载
    lazyload＝false 
    ;是否开始监控所有roots目录下的配置文件变化，改变自动处理(增加，删除，修改)
    watch=true
    ;global，全局属性 是否启用查询数据缓存功能，启用则配置sql中配置属性cached=true 才有效，否则一律不缓存 
    cached=true
    ;global 全局属性 默认缓存的时间，如果使用缓存并且未配置缓存时间则使用该默认时间，单位为分钟，-1为一直有效，-2为一月，-3为一周 -4为一天，单位为分钟。
    cachetime=30
    ;SQL配置文件加载的根目录，可以个多个，定义后自动将真实文件名映射到前边名称
    [roots]
    biz=bizmodel  ; 比如： 系统根目录/bizmodel/plan/main.msql 访问url为 bos/biz/plan/main
    sys=sysmodel  ; 比如： 系统根目录/sysmodel/admin/users.msql 访问url为 bos/sys/admin/users `
## SQL节点属性
    <sql type=""  id=""  eachtran="false/true">
       type=全局SQL类型
       id= 此SQSL唯一标识（调用依据）
       desc=此SQL说明
       eachtran= 此节点（batchexec、batchmultiexec有效，其他类型无效）生成的所有SQL在一个事务中执行，默认为 false 
                       true 则每个批次的SQL在一个事务中执行。     
 
## 全部SQL类型
    配置类型                      内部类型                          说明 
    - select                       ST_SELECT       Tsqltype = iota  //0=普通查询,只能配置一个查询语句返回一个结果集   
    - pagingselect                 ST_PAGINGSELECT                  //1=分页查询，配置一个总数sql语句，一个分页查询语句，返回总页数与一个结果集
    - multiselect                  ST_MULTISELECT                   //3=多结果集查询，配置多个查询语句并返回多个结果集
    - exec/delete                  ST_EXEC                          //4=删除，可以配置多个sql语句，在一个事务中依次执行 
           insert                  ST_EXEC                          //  插入，可以配置多个sql语句，在一个事务中依次执行
           update                  ST_EXEC                          //  更新，可以配置多个sql语句，在一个事务中依次执行 
    - batchexec/batchinsert        ST_BATCHEXEC                     // 单个事务中执行 5=批量插入，配置多个sql，根据参数循环执行(单数据集批量插入)
                batchupdate        ST_BATCHEXEC                     //  批量更新，配置多个sql，根据参数循环执行(单数据集批量更新)
    - batchmultiexec/batchcomplex  ST_BATCHMULTIEXEC                //6=批量复合SQL，配置多个sql(一般多数据集批量插入或更新)
    - getblob                      ST_GETBLOB                       //从数据库获取二进制内容
    - setblob                      ST_SETBLOB                       //保存二进制内容到数据库

## 客户端传入参数
    - select/pagingselect/multiselect/exec(delete/insert/update/getblob/setblob)参数,简单json参数 
         {"id":"001","name":"aaaaa"}
        其中参数中可包含：
         - "callback":"可选参数，不为空时返回JSONP"（仅适用于 select/pagingselect/multiselect）,
         - "start":"可选参数，分页查询时需要 开始记录数"（仅适用于 pagingselect）,
         - "limted":"可选参数，分页查询时需要 每页记录数"（仅适用于 pagingselect）,
  

    - batchexec(batchinsert/batchupdate)参数,简单批量json参数(数组)
        [{"id":"001","name":"aaaaa"},{"id":"002","name":"bbbbb"},{"id":"003","name":"ccc"}]
        分别循环数组中的 json对象调用配置的sql
    -  
    - batchmultiexec/batchcomplexc参数,复杂批量json参数
        {"main":[{"id":"01","name":"aaaaa"},{"id":"002","name":"bbbbb"}],
        "sub1":[{"id":"0111","pid":"01","name":"sub11"},{"id":"0112","pid":"01","name":"sub12"}]
        "sub2":[{"id":"0121","pid":"01","name":"sub21"},{"id":"0122","pid":"01","name":"sub22"}]
        }
         参数main,sub1,sub2分别对应配置的三个cmd的sql参数(每个cmd可以批量执行，同batchexec)

## 参数配置
    `<parameter type="string" name="code" desc="用户帐号" required="true" minlen="5" maxlen="50" default="usercode" return="true"/> `
    - parameter的属性说明
         - name: 参数名称，必须与SQL中的参数名称对应，即?name
         - desc: 参数描述
         - type:参数类型 string/int/float/date/datetime/email/blob 会根据类型进行验证                 
         - required: true 不能为空，也就是其转换为字符串长度>0   false 可以为空
         - minlen: 最小长度（只适用于 type=string） 
         - maxlen: 最大长度（只适用于 type=string） 
         - minvalue: 最小值（只适用于 type=int/float) 
         - maxvalue: 最大值（只适用于 type=int/float) 
         - cached : 是否缓存结果，0=不缓存 1=缓存，缓存的时间由cachetime确定（如果没有配置cachetime则自动为30分钟），只对 select，multiselect，pagingselect(第一页)有效
         - cachetime ：缓存有效时间，不配置或配置为0时 默认为directsql.config的参数分钟，-1为一直有效，-2为一月，-3为一周，单位为分钟。 
         - return: 是否返回，0或不配置为不返回，1为返回该值到客户端，(只适用于带有服务端默认值的才起作用)
         - parentid 是否是作为parentid 使用，0或不配置则不作为父id使用，配置为1则作为从表的与主表关联的父id使用，在SQL类型为batchcomplexc 的作为主从表（一主多从，主表只有一条记录）的从表的父id使用，从表的 SQL参数中需要配置 default类的取值为 parentid，则系统自动用主表的这个值设置到从表的这个参数值中  
         - default: 服务端默认值：如果存在服务端默认值定义则客户端传入参数时可以不传由服务端处理并不执行验证规则（如果客户端传入了则使用客户端的值，并执行服务端的规则验证），
            默认参数取值如下：
             - uuid: 生成新的uuid(36位全球唯一id)
             - int64uuid ：64位整数长度的唯一id 
             - shortuuid ：短字符串的唯一id（将int64uuid转为36进制的值（10个数字+26个字母组成的））
             - nowdate: 当前服务器日期
             - now: 当前服务器日期时间 
             - nowunix: 当前服务器日期时间 unix格式(64位整数)
             - parentid :父id，在SQL类型为batchcomplexc 的作为主从表（一主多从，主表只有一条记录）的从表的父id使用
             - {value} ：直接默认值，用{...}包裹的值直接作为参数默认值,例如 { },{0}
             - 自注册函数取值默认值                
                默认参数取值扩展定义并使用
               1）编写函数 参数必须为*faygo.Context:
                  func name(ctx *faygo.Context) interface{
                      ...
                  }
               2）注册函数
                  RegAny("name",func)
               3）SQL默认参数配置，系统会自动解析调用
                <parameters>
                    <parameter name="id" desc="用户id" type="string" required="true" maxlen="36" minlen="36" default="name" return="true"/>
                </parameters>

## 查询结果缓存
    - 首先在 directsql.ini中进行配置全局参数
       ;是否启用查询数据缓存功能，启用则配置sql中配置属性cached=true 才有效，否则一律不缓存
       cached=true
       ;默认缓存的时间，如果使用缓存并且未配置缓存时间则使用该默认时间，单位为分钟，-1为一直有效，-2为一月，-3为一周 -4为一天，单位为分钟。
      defaultcachetime=30    
    - 在sql配置文件中的sql节点配置 属性
       - cached : 是否缓存结果，0=不缓存 1=缓存，缓存的时间由cachetime确定（如果没有配置cachetime则自动为30分钟），只对 select，multiselect，pagingselect(第一页)有效
       - cachetime ：缓存有效时间，不配置或配置为0时 默认为directsql.ini的参数分钟，-1为一直有效，-2为一月，-3为一周，单位为分钟。 
    - 说明
       - 缓存的key值用 执行请求的路径(/sys/home/select)， 参数名与参数值对作为suffix，进行确定换成值，对于同一个sql只缓存一次，就是第一次执行的参数的结果，其他的参数查询不缓存；对于没有参数的结果缓存suffix=nil  

## 完整示例
    ```<!-- id为本model的标识一般同文件名，database为xorm.ini中配置的数据库名称，为执行该配置文件sql的连接，空为默认数据库 -->
       <model id="demo" database="">
        <comment>
            <desc>DirectSQL功能测试SQL定义</desc>
            <author></author>
            <date>2016-06-20</date>
            <history Editor="畅雨" Date="2016-08-03">完善并增加注释</history>
        </comment>
        <!-- JSON参数示例：{'code':'12345'} -->
        <sql type="select" id="select" desc="查询SQL,结果缓存30分钟" cached="true" cachetime="30">
            <cmd in="" out="">
                <![CDATA[ select * from sys_user where code=?code ]]>
                <parameters>
                    <parameter type="string" name="code" desc="用户帐号" required="true" minlen="5" maxlen="50"/>
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：{'name':'李四','start':1,'limted':20 } -->
        <sql type="pagingselect" id="paging" desc="服务端分页查询,第一个cmd=总页数SQL,第二个cmd=分页的查询数据sql">
            <cmd in="" out="total">
                <![CDATA[ SELECT count(id) AS total FROM sys_user ]]>
            </cmd>
            <cmd in="" out="data">
                <![CDATA[ SELECT * FROM sys_user LIMIT ?start,?limted ]]>
                <parameters>
                    <parameter  name="start" desc="start" type="int" required="true" />
                    <parameter  name="limted" desc="size" type="int" required="true" />
                </parameters>
            </cmd>
        </sql>
         <!-- JSON参数示例：{'code':'12345'} -->
        <sql type="nestedselect" id="nested" desc="嵌套的结果json---已作废" idfield="id" pidfield="pid">
            <cmd in="" out="">
                <![CDATA[ SELECT * FROM sys_sys_user ]]>
            </cmd>
        </sql>
         <!-- JSON参数示例：{'code':'12345','nick':'changyu'} -->
        <sql type="multiselect" id="multi" desc="多个Select返回多个结果集的查询组合，一个json参数">
            <cmd in="" out="main">
                <![CDATA[  SELECT id,code FROM sys_user where code=?code ]]>
                <parameters>
                    <parameter type="string" name="code" desc="用户帐号" required="true" minlen="5" maxlen="50" />
                </parameters>
            </cmd>
            <cmd in="" out="detail1">
                <![CDATA[   SELECT id,code,pwd,nick FROM sys_user   ]]>
            </cmd>
            <cmd in="" out="detail2">
                <![CDATA[   SELECT id,code,cnname,pwd,nick FROM sys_user WHERE nick=?nick   ]]>
                <parameters>
                    <parameter type="string" name="nick" desc="用户帐号" required="true" />
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：{'code':'123'} -->
        <sql type="delete" id="delete" desc="删除，可以执行多条删除语句,建议类型改为 exec">
            <cmd in="" out="">
                <![CDATA[   DELETE FROM sys_user where code=?code   ]]>
                <parameters>
                    <parameter name="code" desc="用户帐号" type="string" required="true" />
                </parameters>
            </cmd>
            <cmd in="" out="">
                <![CDATA[   DELETE FROM sys_userdetail where pcode=?code   ]]>
                <parameters>
                    <parameter name="code" desc="用户帐号" type="string" required="true" />
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：{'code':'123','name':'xxx'} -->
        <sql type="insert" id="insert" desc="新增服务端生成newid返回,建议类型改为 exec">
            <cmd in="" out="">
                <![CDATA[       INSERT INTO sys_user(id,code,cnname,pwd,nick) VALUES(?id,?code,?cnname,?pwd,?nick) ]]>
                <parameters>
                    <parameter name="id" desc="用户id" type="string" required="true" maxlen="36" minlen="36" default="uuid" return="true"/>
                    <parameter name="code" desc="用户帐号" type="string" required="true" maxlen="50" minlen="5" />
                    <parameter name="cnname" desc="cnname" type="string" required="true" />
                    <parameter name="nick" desc="昵称" type="string" required="true" maxlen="30" minlen="0" />
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：{'code':'123','nick':'xxx'} -->
        <sql type="update" id="update" desc="更新,建议类型改为 exec" >
            <cmd in="" out="">
                <![CDATA[ update sys_user set nick=?nick where code=?code ]]>
                <parameters>
                    <parameter type="string" name="code" desc="用户帐号" required="true" />
                    <parameter type="string" name="nick" desc="用户帐号" required="true" />
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：{'code':'123','nick':'xxx'} -->
        <sql type="insert" id="save" desc="保存(插入或更新),建议类型改为 exec" >
            <cmd in="" out="">
                <![CDATA[ INSERT INTO sys_user(uid,code,pwd,nick,login_ip,login_time,regist_time) VALUES(?uid,?code,?pwd,?nick,?login_ip,?login_time,?regist_time) ON DUPLICATE KEY UPDATE nick=?nick   ]]>
                <parameters>
                    <parameter type="string" name="id" desc="id" required="true" default="uuid" return="true"/>
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：  [{"id":"001","name":"aaaaa"},{"id":"002","name":"bbbbb"},{"id":"003","name":"ccc"}] -->
        <sql type="batchinsert" id="batchinsert" desc="批量新增，json传来多条数据，一个批次插入，要么全部成功要么全部失败">
            <cmd in="" out="">
                <![CDATA[ INSERT INTO sys_user(id,code,cnname,pwd,nick) VALUES(?id,?code,?cnname,?pwd,?nick)            ]]>
                <parameters>
                    <parameter type="string" name="id" desc="用户帐号" required="true" minlen="30" default="uuid" return="true"/>
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：[{"id":"001","name":"aaaaa"},{"id":"002","name":"bbbbb"},{"id":"003","name":"ccc"}]-->
        <sql type="batchupdate" id="bacthsave" desc="批量更新,json传来多条数据，一个批次更新，要么全部成功要么全部失败">
            <cmd in="" out="">
                <![CDATA[   INSERT INTO sys_user(id,code,cnname,pwd,nick) VALUES(?id,?code,?cnname,?pwd,?nick)     ]]>
            </cmd>
        </sql>
            <!-- JSON参数示例：  [{"id":"001","name":"aaaaa"},{"id":"002","name":"bbbbb"},{"id":"003","name":"ccc"}]  -->
        <sql type="batchdelete" id="delete" desc="批量删除，根据参数多次执行该语句">
            <cmd in="" out="">
                <![CDATA[   DELETE FROM sys_user where code=?code   ]]>
                <parameters>
                    <parameter name="code" desc="用户帐号" type="string" required="true" />
                </parameters>
            </cmd>
        </sql>
        <!-- JSON参数示例：{main:[{'id':'123','name':'xxx'}],detail1:[{},...],detail2:[{},...]} -->
        <sql type="batchmultiexec" id="bacthcomplexsave" desc="批量在一个事务里边组合执行多个表的保存，要么全部成功要么全部失败">
            <cmd in="main" desc="主表的数据，支持多条">
                <![CDATA[   INSERT INTO sys_user(id,code,cnname,pwd,nick) VALUES(?id,?code,?cnname,?pwd,?nick) ]]>
                <parameters>
                    <parameter type="string" name="id" desc="用户帐号" required="true" default="uuid" />
                </parameters>
            </cmd>
            <cmd in="detail1" desc="子表一数据，支持多条">
                <![CDATA[  INSERT INTO sys_user(id,code,cnname,pwd,nick) VALUES(?id,?code,?cnname,?pwd,?nick)]]>
            </cmd>
            <cmd in="detail2" desc="子表二数据，支持多条">
                <![CDATA[  INSERT INTO sys_user(id,code,cnname,pwd,nick) VALUES(?id,?code,?cnname,?pwd,?nick)   ]]>
            </cmd>
            <!-- 下面可继续添加多个组合进来的sql -->
        </sql>
        <!-- 下面可继续添加新的方法 -->
      </model>
   ```

## 存在问题
    - nestedselect-返回嵌套JSON的未实现，放到客户端用js实现
      参考函数：
      `
            /**
             * 将通过id，parentid表达的父子关系列表转换为树形
             * @param {object} data 数据列表
             * @param {string} id 主键
             * @param {string} parentid 父id
             * @returns {json} 
             */
             getNestedJSON = function(data, id, parentid) {
                if ((data === null) || (data == "[]")) {
                    return [];
                }
                // create a name: node map
                var dataMap = data.reduce(function(map, node) {
                    map[node[id]] = node;
                    return map;
                }, {});
                // create the tree array
                var tree = [];
                data.forEach(function(node) {
                    // add to parent
                    var parent = dataMap[node[parentid]];
                    if (parent) {
                        // create child array if it doesn't exist
                        (parent.children || (parent.children = []))
                        // add node to child array
                        .push(node);
                    } else {
                        // parent is null or missing
                        tree.push(node);
                    }
                });
                return tree;
            }
      `
    - sqlwatcher中不能监控改名的文件（可以通过重新载入解决）