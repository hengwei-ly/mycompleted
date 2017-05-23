/**
 *
 * Created by pc on 2017/4/5.
 */

$(function(){
    //设置两个多选框
    getAndWriteAddress();
    $("#filter_facility_selection").select2();
    //设置servity单击事件
    $('.filter_severity_range').click(function () {
        var range_text = $(this).text();
        $('#filter_severity_range_text').text(range_text);
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
    //表单提交事件
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
                    window.location.href = "index.html";
                }else {
                    showErrors(XmlHttpRequest);
                    return;
                }
            }
        });
    });

});

//将错误信息显示在页面上
function showMessages(messages) {
    var templ = document.getElementById('filter_error_messages').innerHTML;
    var errors = {messages:messages};
    var html = ejs.render(templ,errors);
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

//显示来自服务端的错误
function showErrors(err) {
    var templ = document.getElementById('errors_templ').innerHTML;
    var data = {errors:err };
    var html = ejs.render(templ, data);
    $('#syslog_error_brother+div').remove();
    $('#syslog_error_brother').after(html);
}

//获取来自服务器的所有address 并加入到select中
function getAndWriteAddress() {
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

            $("#filter_address_selection").select2();
        }
    });

}

