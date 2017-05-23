/**
 *
 * Created by pc on 2017/4/5.
 */



$(function(){
    //显示原有信息
    var params = window.location.search;
    var id = params.replace("?id=","");

    //设置servity单击事件
    $('.filter_severity_range').click(function () {
        var range_text = $(this).text();
        $('#filter_severity_range_text').text(range_text);
    });
    $.ajax({
        url:"/filters/"+id,
        dataType:"json",
        contentType:"application/json",
        error:function (XmlHttpRequest) {
            showErrors(XmlHttpRequest);
            return;
        },
        success:function (query) {
            //将读取的参数写入到表单中
            WriteWebParams(query);
            //添加验证信息
            $('#filters_update_form').validate({
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
            $('#filters_update_form').submit(function (event) {
                event.preventDefault();
                //将表单中的错误信息写到顶部错误框
                var errors = $('[class="error"]:visible');
                if (errors.length){
                    var messages = new Array();
                    for(var i=0;i<errors.length;i++){
                        var labelname = $(errors[i]).parent().prev().text();
                        var message = $(errors[i]).text();
                        messages[i] = labelname  + message;
                    }
                    showMessages(messages);
                    return;
                }
                //验证无误
                var query = ReadFilterParams();
                query.id = id;
                var param = JSON.stringify(query);
                $.ajax({
                    url:"/filters/"+id,
                    type: "PUT",
                    data:param,
                    dataType : "json",
                    contentType: "application/json",
                    complete:function (XmlHttpRequest) {
                        if (XmlHttpRequest.responseText=="OK"){
                            window.location.href = "index.html";
                        }else {
                            showErrors(XmlHttpRequest);
                            return;
                        }
                    }
                });
            });
        }
    });
});

//获取来自服务器的所有address 并加入到select中
function getAndWriteAddress(values) {
    $.ajax({
        url:"/fields/address",
        type:"GET",
        dataType:"json",
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
            $('#filter_address_selection').children().remove();
            $('#filter_address_selection').append(html);

            $('#filter_address_selection').val(values);
            $("#filter_address_selection").select2();
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

//显示来自服务端的错误
function showErrors(err) {
    var templ = document.getElementById('errors_templ').innerHTML;
    var data = {errors:err };
    var html = ejs.render(templ, data);
    $('#syslog_error_brother+div').remove();
    $('#syslog_error_brother').after(html);
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
        filters[fn] = {
            field:"facility",
            op:"QueryString",
            values:facility};
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

//写入原本数据
function WriteWebParams(query) {
    $('#filters_name').val(query.name);
    $('#filters_description').val(query.description);
    console.log(query);
    if (query.filters){
        var address = new Array();
        var severity = new Array();
        var facility = new Array();
        var message = new Array();
        for (var i=0;i<query.filters.length;i++){
            var filter = query.filters[i];
            if (filter.field=="address"){
                address = filter.values;
                continue;
            }
            if (filter.field=="severity"){
                severity = filter.values;
                continue;
            }
            if (filter.field=="facility"){
                facility = filter.values;
                continue;
            }
            if (filter.field=="message"){
                message = filter.values;
                continue;
            }
        }
        //写入address
        getAndWriteAddress(address);
        //severity
        var start_severity = severity[0];
        var end_severity = severity[1];
        if (start_severity == end_severity) {
            $('#filter_severity_range_text').text("等于");
            $('#filter_severity_selection').val(start_severity);
        }else if (end_severity == "7") {
            $('#filter_severity_range_text').text("大于等于");
            $('#filter_severity_selection').val(start_severity);
        }else if (start_severity=="0") {
            $('#filter_severity_range_text').text("小于等于");
            $('#filter_severity_selection').val(end_severity);
        }

        $("#filter_facility_selection").val(facility);
        $("#filter_facility_selection").select2();

        $('#filter_message').val(message[0]);
    }
}

