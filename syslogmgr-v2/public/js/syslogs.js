/**
 * Created by pc on 2017/3/29.
 * 更改条件后重新取值，如果恢复默认值，直接写默认值，如果要保持不变，则从对应的控件中取值
 * 所有请求参数封装在urlParams里面  如果需要再进行增加条件：
 *  ReadParams  WriteParams 直接添加，然后再添加对对应参数监听的方法即可
 *
 */
$(function(){
    //创建filter下拉选
    CreateFilterUnits();
    //创建日期选择框
    CreateDateSelector();
    //创建临时查询选择框
    // createFastFilter();
    //改变显示页数
    $('#syslog_pagesize_select').change(function () {
        var params = ReadParams();
        if (params.filterId=="fast"){
            fastQueryForSyslogs(params);
        }else {
            ListSyslogs(params);
        }
        return;
    });

    //更改过滤条件
    $('#syslog_filter_select').change(function () {
        var params = ReadParams();
        if (params.filterId=="fast"){
            $('#fast_filters_form').attr('hidden',false);
            $('#fast_filters_form').find("*").removeAttr('hidden');
            //创建一个临时表单用于选择过滤条件
            createFastFilter();
            //隐藏前两项  name & description
            $('#fast_filters_form').find('.form-group:lt(3)').attr('hidden',true);
        }else {
            $('#fast_filters_form').attr('hidden',true);
            ListSyslogs(params);
        }
        return;
    });

    $('#fast_new_filter').click(function () {
    });

    //导出为csv
    $('#syslog_out_csv').click(function () {
        getLoadSources(getSyslog4CSV);
    });
    //导出为txt
    $('#syslog_out_txt').click(function () {
        getLoadSources(getSyslog4TXT);
    });

    var urlParams = ReadParams();
    ListSyslogs(urlParams);
});

//写入参数
function WriteParams(urlParams){
    //清空错误信息
    $('#syslog_errors').children().remove();
    //清空表格内数据
    $('#datas_body tr').remove();
    //移除页面选择按钮
    $('#pages_body ul').remove();
    //初始化记录数目显示
    $('#total_counts').html("0");
    $('#syslog_pagesize_select').val(urlParams.pagesize);
    $('#currentPage').val(urlParams.page);
    $('#syslog_filter_select').val(urlParams.filterId);
    $('#starttime').val(urlParams.starttime);
    $('#endtime').val(urlParams.endtime);
}

//读取参数
function ReadParams(){
    var pagesize = $('#syslog_pagesize_select').val();
    var filterId = $('#syslog_filter_select').val();
    var starttime = $('#starttime').val();
    var endtime = $('#endtime').val();
    var urlParams = {
        page:1,
        pagesize:pagesize,
        filterId:filterId,
        starttime:starttime,
        endtime:endtime
    };
    return urlParams;
}

//列出数据
function ListSyslogs(urlParams) {
    //清空原纪录并记录当前查询值
    WriteParams(urlParams);
    var templ = document.getElementById('syslogs_table').innerHTML;
    //最大页数
    var maxpage = 0;
    var countUrl = "/query/"+urlParams.filterId+"/count?start_at="+urlParams.starttime+"&end_at="+urlParams.endtime;
    $.ajax({
        url:countUrl,
        type: "GET",
        error:function (XmlHttpRequest) {
            showErrors(XmlHttpRequest);
            return;
        },
        success:function (count) {
            if (count==NaN){
                showErrors(XMLHttpRequest(this));
                return;
            }
            $('#total_counts').text(count);
            maxpage = Math.ceil(count/urlParams.pagesize);
            if (maxpage==NaN || maxpage<=0){
                maxpage = 1;
            }
            //根据最大页数创建选页框
            CreatePageUnits(maxpage);
            var queryURL ="/query/"+urlParams.filterId+"?offset="+urlParams.pagesize*(urlParams.page-1)+
                        "&limit="+urlParams.pagesize+"&start_at="+urlParams.starttime+"&end_at="+urlParams.endtime;
            // console.log(queryURL);
            $.ajax({
                url:queryURL,
                type: "GET",
                dataType : "json",
                contentType: "application/json",
                error:function (XmlHttpRequest) {
                    showErrors(XmlHttpRequest);
                    return;
                },
                success:function (datas) {
                    var data = {syslogs: datas};
                    var html = ejs.render(templ, data);
                    //清空模板然后插入模板
                    $('#datas_body tr').remove();
                    $('#datas_body').prepend(html);
                    $('#aim_page').val(urlParams.page);
                    //改为文字
                    toChword();
                    //渲染颜色
                    addColor();
                    for (var i=0;i<2;i++){
                        $('#date_select1').trigger('click');
                    }
                }
            });
        }
    });
}

//创建页面选择按钮
function CreatePageUnits(maxpage) {
    var templ = document.getElementById('pageUnits').innerHTML;
    var currentPage = $('#currentPage').val();
    if (currentPage==NaN || currentPage<=0){
        currentPage = 1;
    }else {
        currentPage = ($('#currentPage').val())*1;
    }
    var pages = {maxpage:maxpage,currentPage:currentPage};
    var html = ejs.render(templ,pages);
    $('#pages_body ul').remove();
    $('#pages_body').prepend(html);
}

//跳转页面功能
function ToPage(pageValue){
    var pagesize = ($('#syslog_pagesize_select').val())*1;
    var currentPage = ($('#currentPage').val())*1;
    var maxPage = ($('#maxPage').text())*1;

    var page = $(pageValue).text();
    if (page=='跳转'){
        page = $('#aim_page').val();
        if(page==NaN){
            return
        }
        page = ($('#aim_page').val())*1;
        if (page<1) {
            page = 1;
        }else if (page> maxPage) {
            page = maxPage;
        }
    }else if (page=='上一页'){
        if (currentPage==1) {
            return
        }
        page = currentPage-1;
    }else if (page=='下一页'){
        if (currentPage>=maxPage) {
            return
        }
        page = currentPage+1;
    }else if (page=='首页'){
        page = 1;
    }else if (page=='尾页'){
        page = maxPage;
    }else{
        page = eval(page);
        if (page<1 || page>maxPage){
            page = 1;
        }
    }
    //读取原数据
    var params = ReadParams();
    params.page = page;
    params.pagesize = pagesize;
    if (params.filterId=="fast"){
        fastQueryForSyslogs(params);
    }else {
        ListSyslogs(params);
    }
}

//按数字等级转为文字 设备号转换为汉字
function toChword() {
    var docs = $('[type="severities"]');
    for(var i=0;i<docs.length;i++){
        var level = $(docs[i]).text();
        if(level=="0"){
            $(docs[i]).text("Emergency");
        }else if(level=="1"){
            $(docs[i]).text("Alert");
        }else if(level=="2"){
            $(docs[i]).text("Critical");
        }else if(level=="3"){
            $(docs[i]).text("Error");
        }else if(level=="4"){
            $(docs[i]).text("Warning");
        }else if(level=="5"){
            $(docs[i]).text("Notice");
        }else if(level=="6"){
            $(docs[i]).text("Informational");
        }else if(level=="7"){
            $(docs[i]).text("Debug");
        }else {
            $(docs[i]).text("其他等级");
        }
    }

    var fs = $('[type="facilities"]');
    for (var i=0;i<fs.length;i++){
        var fac = $(fs[i]).text();
        if (fac=="0"){
            $(fs[i]).text("内核");
        }else if(fac=="1"){
            $(fs[i]).text("用户");
        }else if(fac=="2"){
            $(fs[i]).text("邮件");
        }else if(fac=="3"){
            $(fs[i]).text("系统进程");
        }else if(fac=="4"){
            $(fs[i]).text("安全认证");
        }else if(fac=="5"){
            $(fs[i]).text("syslogd内部");
        }else if(fac=="6"){
            $(fs[i]).text("打印系统");
        }else if(fac=="7"){
            $(fs[i]).text("网络");
        }else if(fac=="8"){
            $(fs[i]).text("UUIP");
        }else if(fac=="9"){
            $(fs[i]).text("时钟");
        }else if(fac=="10"){
            $(fs[i]).text("安全授权");
        }else if(fac=="11"){
            $(fs[i]).text("FTP");
        }else if(fac=="12"){
            $(fs[i]).text("NTP");
        }else if(fac=="13"){
            $(fs[i]).text("日志审计");
        }else if(fac=="14"){
            $(fs[i]).text("日志警报");
        }else if(fac=="15"){
            $(fs[i]).text("调度程序");
        }else if(fac=="16"){
            $(fs[i]).text("本地0");
        }else if(fac=="17"){
            $(fs[i]).text("本地1");
        }else if(fac=="18"){
            $(fs[i]).text("本地2");
        }else if(fac=="19"){
            $(fs[i]).text("本地3");
        }else if(fac=="20"){
            $(fs[i]).text("本地4");
        }else if(fac=="21"){
            $(fs[i]).text("本地5");
        }else if(fac=="22"){
            $(fs[i]).text("本地6");
        }else if(fac=="23"){
            $(fs[i]).text("本地7");
        }else {
            $(fs[i]).text("其他设备");
        }
    }

}

//按等级标记颜色
function addColor() {
    var level_colors = {
        "Emergency":"#ff0000",
        "Alert":"#dd2222",
        "Critical":"#ed5565",
        "Error":"#C274BB",
        "Warning":"#FF8800",
        "Notice":"#1111ff",
        "Informational":"#66c0b0",
        "Debug":"#22ff33"
    };
    var servrities = document.getElementsByClassName("label label-primary");
    for (var i=0;i<servrities.length;i++){
        var color_index = $(servrities[i]).text().trim();
        $(servrities[i]).css("background-color",level_colors[color_index]);
    }
}

//获取filterNames 和 FilterIds并写入html中
function CreateFilterUnits () {
    $.ajax({
        url:"/filters",
        dataType:"json",
        async:false,
        error:function (XmlHttpRequest) {
            var message = "过滤条件加载错误";
            showErrors(XmlHttpRequest,message);
            return;
        },
        success:function (querys){
            if (!querys){
                return;
            }
            var templ = document.getElementById('filters_body').innerHTML;
            var filters = {querys:querys};
            var html = ejs.render(templ,filters);
            $('#syslog_filter_select option:gt(1)').remove();
            $('#syslog_filter_select').append(html);
            if ($('#syslog_filter_select option').length<=1){
                return;
            }
            $('#syslog_filter_select option:eq(0)').attr("selected",true);
            return;
        }
    });
}

//创建日期选择器并进行请求
function CreateDateSelector() {
    //初始化选择器的起止时间
    var endDate = new Date();
    var e_year = endDate.getFullYear();
    var e_month = endDate.getMonth()+1;
    var e_day = endDate.getDate();
    //默认显示前三天 可修改数值
    var preday = 3;
    var startDate = new Date(endDate.getTime()-1000*60*60*24*preday);
    var s_year = startDate.getFullYear();
    var s_month = startDate.getMonth()+1;
    var s_day = startDate.getDate();

    var dateTempl = document.getElementById('dateTempl').innerHTML;

    var starttime = s_year + "/" + s_month + "/" + s_day;
    var endtime = e_year + "/" + e_month + "/" + e_day;
    var datas = {date:{starttime:starttime,endtime:endtime}};
    var html = ejs.render(dateTempl,datas);
    $('#date_select1').popover({
        trigger:"click",
        placement:"auto bottom",
        html:true,
        content:html
    });
    //创建日期选择弹框
    createDatePicker();
    var n = 0;
    $(document).click(function (event) {
        if (n<3){
            return
        }
        var popattr = $('#date_select1').attr("aria-describedby");
        var popval = $('#popoverid').val();
        if (popattr == popval || popval =="popover"){
        }else{
            //如果弹出的框存在
            if ($('[role="tooltip"]').html()){
                //判断不需要让框消失的控件  如果点击不需要消失的控件  mm=false  空白为true  为true时弹窗消失
                var mm = true;
                var aims = $('[role="tooltip"]').find("*");
                //判断当前弹窗中的控件
                for (var i=0;i<aims.length;i++){
                    if (aims[i]==event.target){
                        mm = false;
                        break;
                    }
                }
                //判断日期弹窗里面的控件控件
                var targetClass = $(event.target).attr('class');
                if (targetClass){
                    if (targetClass.indexOf("year")!=-1){
                        mm = false;
                    }else if (targetClass.indexOf("month")!=-1){
                        mm = false;
                    }else if (targetClass.indexOf("day")!=-1){
                        mm = false;
                    }
                }
                //判断是否为选择的按钮
                if ($(event.target).attr('id')=="date_select1" || $(event.target).parent().attr('id')=="date_select1"){
                    mm = false;
                }
                //如果mm为真  隐藏弹窗 隐藏弹窗起始就两句话
                //1.找到弹窗，清除掉   2.模拟点击一下原控件
                if (mm){
                    $('[role="tooltip"]').remove();
                    $('#date_select1').trigger('click');
                }
            }
        }
    });


    var aims = $('[role="tooltip"]').find("*");
    for (var i=0;i<aims.length;i++){
        if (aims[i]==e){
            return ;
        }
    }
    if ($('.datepicker').length){
        var dates = $('.datepicker').find("*");
        for (var i=0;i<dates.length;i++){
            if (dates[i]==e){
                return ;
            }
        }
    }

    $('#date_select1').click(function () {
        var popid = $('#date_select1').attr("aria-describedby");
        if (n%2==1){
            $('#popoverid').val(popid);
        }
        n++;
    });
}

//创建日期选择弹框 和 选择后提交
function createDatePicker() {
    $("#date_select1").on('shown.bs.popover', function () {
        $('#data_5 .input-daterange').datepicker({
            keyboardNavigation: false,
            forceParse: false,
            autoclose: true,
            format: "yyyy/mm/dd"
        });

        //更改起止时间
        $('#select_rangedate_confirm').click(function () {
            var params = ReadParams();
            var start = $('#select_date_start').val();
            var end = $('#select_date_end').val();
            params.starttime = new Date(start).toISOString();
            //結束時間要加一天
            params.endtime = new Date(new Date(end).getTime()+60*60*24*1000-1).toISOString();
            $("#date_select1").trigger("click");
            $('#date_select1 span').text(start+"至"+end);

            if (params.filterId=="fast"){
                fastQueryForSyslogs(params);
            }else {
                ListSyslogs(params);
            }
            return;
        });

        //选择时间
        $('#condate_selectors button').each(function () {
            $(this).click(function () {
                var dateId = $(this).attr("date");
                var params = ReadParams();
                params.endtime = new Date().toISOString();
                //通过标签获取需要查询的值
                params.starttime = exExsitsDate(dateId).toISOString();
                $("#date_select1").trigger("click");
                $('#date_select1 span').text($(this).text());

                if (params.filterId=="fast"){
                    fastQueryForSyslogs(params);
                }else {
                    ListSyslogs(params);
                }
                return;
            });

        });

    });
}

//根据指定标签  返回起始日期
function exExsitsDate(dateId) {
    var date = new Date();
    if (dateId=="one_day"){
        var times = 24*60*60*1000*1;
        return new Date(date.getTime()-times);
    }
    if (dateId=="three_days"){
        var times = 24*60*60*1000*3;
        return new Date(date.getTime()-times);
    }
    if (dateId=="one_week"){
        var times = 24*60*60*1000*7;
        return new Date(date.getTime()-times);
    }
    if (dateId=="one_month"){
        var times = 24*60*60*1000*30;
        return new Date(date.getTime()-times);
    }
    if (dateId=="three_months"){

        return new Date(date.getFullYear(),date.getMonth()-3,date.getDate());
    }
    if (dateId=="one_year"){
        return new Date(date.getFullYear()-1,date.getMonth(),date.getDate());
    }
}

//显示错误信息
function showErrors(err) {
    var templ = document.getElementById('errors_templ').innerHTML;
    var data = {errors:err};
    var html = ejs.render(templ, data);
    $('#syslog_errors').children().remove();
    $('#syslog_errors').append(html);
}

//获取需要下载的数据的url
function getLoadSources(action) {
    var urlParams = ReadParams();
    urlParams.filterId = urlParams.filterId=="fast" ? "0":urlParams.filterId;
    //最大查询条目数
    urlParams.pagesize = "10000000";
    var queryURL ="/query/"+urlParams.filterId+"?offset="+urlParams.pagesize*(urlParams.page-1)+
        "&limit="+urlParams.pagesize+"&start_at="+urlParams.starttime+"&end_at="+urlParams.endtime;
    action(queryURL);
}

//导出为csv文件  创建临时数据表格然后下载
function getSyslog4CSV(queryURL) {
    $.ajax({
        url:queryURL,
        type: "GET",
        dataType : "json",
        contentType: "application/json",
        async:false,
        success:function (datas) {
            var str ="接收时间,设备,等级,主机名,消息";
            for (var i=0;i<datas.length;i++){
                var d = datas[i];
                str += "\n"+d.reception+","+d.priority+","+d.severity+","+d.address+","+d.message;
            }
            str =  encodeURIComponent(str);
            $('#syslog_out_csv').attr('href',"data:text/csv;charset=utf-8,\ufeff"+str);
            $('#syslog_out_csv').attr('download','syslogs.csv');
            return;
        },
        error:function (XmlHttpRequest) {
            showErrors(XmlHttpRequest);
            return;
        }
    });


}

function getSyslog4TXT(queryUrl) {
    $('#syslog_out_txt').attr('href',queryUrl);
    $('#syslog_out_txt').attr('download',"syslogs.txt");
}

//增加临时查询框  并对其中的按钮进行监听
function createFastFilter() {
    //写入address
    getAndWriteAddress();
    $("#filter_facility_selection").select2();
    //设置servity单击事件
    $('.filter_severity_range').click(function () {
        var range_text = $(this).text();
        $('#filter_severity_range_text').text(range_text);
    });

    //查询按钮 --> 按照指定条件直接查询
    $('#filter_once_query').click(function () {
        var errors = $('[class="error"]:visible');
        // 错误消息提示
        var messages = new Array();
        if (errors.length){
            for(var i=0;i<errors.length;i++){
                var labelname = $(errors[i]).parent().prev().text();
                var message = $(errors[i]).text();
                messages[i] = labelname  + message;
            }
            showMessages(messages);
            return;
        }
        //快速匿名查询
        var urlParams = ReadParams();
        fastQueryForSyslogs(urlParams);
    });

    //保存按钮 --> 保存name 和 description
    $('#filter_once_create').click(function () {
        //显示三两项  name & description 和按钮  隐藏后面的中间四项
        $('#fast_filters_form').find('.form-group:lt(3)').attr('hidden',false);
        $('#fast_filters_form').find('.form-group:gt(2)').attr('hidden',true);
    });

    //表单验证
    $('#filters_create_form').validate({
        rules: {
            name: {
                required: true,
                maxlength: 15
            },
            description: {
                required: true
            }
        }
    });

    //创建按钮 --> 表单提交
    $('#filters_create_form').submit(function (event ) {
        event.preventDefault();
        var errors = $('[class="error"]:visible');
        // 错误消息提示
        var messages = new Array();

        //获取表单信息
        var query = ReadFilterParams();

        if (errors.length){
            for(var i=0;i<errors.length;i++){
                var labelname = $(errors[i]).parent().prev().text();
                var message = $(errors[i]).text();
                messages[i] = labelname  + message;
            }
            showMessages(messages);
            return;
        }
        //没有错误信息
        var param = JSON.stringify(query);
        $.ajax({
            type: "POST",
            url: "/filters",
            data:param,
            dataType : "json",
            contentType: "application/json",
            complete:function (XmlHttpRequest) {
                if (XmlHttpRequest.responseText=="OK"){
                    //表单提交成功后，将过滤器选择在当前创建的过滤器上
                    CreateFilterUnits();
                    $('#syslog_filter_select option:eq(0)').attr("selected",false);
                    $('#syslog_filter_select option:last').attr("selected",true);
                    var params = ReadParams();
                    ListSyslogs(params);
                    //确认提交后,将当前临时框隐藏,数据清空(先不清空)
                    $('#fast_filters_form').attr('hidden',true);
                }else {
                    showErrors(XmlHttpRequest);
                    return;
                }
            }
        });
    });

    //取消按钮
    $('#filters_once_concel').click(function (event) {
        event.preventDefault();
        //隐藏前两项  name & description
        $('#fast_filters_form').find('.form-group:gt(2)').attr('hidden',false);
        $('#fast_filters_form').find('.form-group:lt(3)').attr('hidden',true);
    });
}

//读取表单中的参数
function ReadFilterParams() {
    var name = $('#filters_name').val();
    var description = $('#filters_description').val();

    var address = $('#filter_address_selection').val();

    var severity = $('#filter_severity_selection').val();
    var range_text = $('#filter_severity_range_text').text();

    var start_severity = "0";
    var end_severity = "7";
    if (range_text=="等于") {
        start_severity = severity;
        end_severity = severity;
    }else if(range_text=="小于等于"){
        end_severity = severity;
    }else if (range_text = "大于等于"){
        start_severity = severity;
    }

    var facility = $('#filter_facility_selection').val();

    var message = $('#filter_message').val();

    //如果当前项没有选择  对应的filter为空
    var filters = new Array();
    var fn = 0;
    if (address != null){
        filters[fn] = {
            field:"address",
            op:"Term",
            values:address};
        fn++;
    }
    if (facility!=null){
        var value = "";
        for (var i=0;i<facility.length;i++){
            value +=facility[i]+" ";
        }
        filters[fn] = {
            field:"facility",
            op:"QueryString",
            values:[value]};
        fn++;
    }
    if (message!=""){
        filters[fn] = {
            field:"message",
            op:"QueryString",
            values:[message]};
        fn++;
    }
    //将servity加进filters
    filters[fn] = {
        field:"severity",
        op:"NumericRange",
        values:[start_severity,end_severity,"true","true"]};
    fn++;

    var query = {
        name:name,
        description:description,
        filters:filters
    };
    return query;
}

//获取来自服务器的所有address 并加入到select中
function getAndWriteAddress() {
    $.ajax({
        url:"/fields/address",
        type:"GET",
        dataType:"json",
        async:false,
        error:function (XmlHttpRequest) {
            showErrors(XmlHttpRequest);
            return;
        },
        success:function (data) {
            if (!data){
                showMessages("获取主机名错误!");
                return;
            }
            var templ = document.getElementById('filter_address_templ').innerHTML;
            var datas = {anames:data};
            var html = ejs.render(templ,datas);
            $('#filter_address_selection option').remove();
            $('#filter_address_selection').append(html);

            $("#filter_address_selection").select2();
        }
    });

}

//临时查询：传入匿名查询条件  返回列出数据
function fastQueryForSyslogs(urlParams) {
    //清空原纪录并记录当前查询值
    WriteParams(urlParams);
    //查询当前条件
    var query = ReadFilterParams();
    var postParam = JSON.stringify(query);
    //最大页数
    var maxpage =0;
    var filterCountURL =  "/query/count?start_at="+urlParams.starttime+"&end_at="+urlParams.endtime;
    var templ = document.getElementById('syslogs_table').innerHTML;
    $.ajax({
        type: "POST",
        url: filterCountURL,
        data:postParam,
        contentType: "application/json",
        error:function (XmlHttpRequest) {
            showErrors(XmlHttpRequest);
            return;
        },
        success:function (count) {
            if (count==NaN){
                showErrors(XMLHttpRequest(this));
                return;
            }
            $('#total_counts').text(count);
            maxpage = Math.ceil(count/urlParams.pagesize);
            if (maxpage==NaN || maxpage<=0){
                maxpage = 1;
            }
            //根据最大页数创建选页框
            CreatePageUnits(maxpage);

            var filterqueryURL = "/query?offset="+urlParams.pagesize*(urlParams.page-1)+
                "&limit="+urlParams.pagesize+"&start_at="+urlParams.starttime+"&end_at="+urlParams.endtime;
            //post方式进行快速查询
            $.ajax({
                url:filterqueryURL,
                type: "POST",
                data:postParam,
                dataType : "json",
                contentType: "application/json",
                error:function (XmlHttpRequest) {
                    showErrors(XmlHttpRequest);
                    return;
                },
                success:function (datas) {

                    var data = {syslogs: datas};
                    var html = ejs.render(templ, data);
                    //清空模板然后插入模板
                    $('#datas_body tr').remove();
                    $('#datas_body').prepend(html);
                    $('#aim_page').val(urlParams.page);
                    //改为文字
                    toChword();
                    //渲染颜色
                    addColor();
                    for (var i=0;i<2;i++){
                        $('#date_select1').trigger('click');
                    }
                    // //查询完成  将临时框隐藏
                    // $('#fast_filters_form').attr('hidden','true');
                    return
                }
            });
        }
    });
}

//将错误信息显示在页面上
function showMessages(messages) {
    var templ = document.getElementById('filter_error_messages').innerHTML;
    var errors = {messages:messages};
    var html = ejs.render(templ,errors);
    $('#syslog_error_brother+div').remove();
    $('#syslog_error_brother').after(html);
}
