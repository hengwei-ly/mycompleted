/**
 *
 * Created by pc on 2017/4/5.
 */

$(function(){
    $.ajax({
        url:"/filters",
        dataType:"json",
        success:function (querys){
            var templ = document.getElementById('querys_body').innerHTML;
            var filters = {querys:querys};
            var html = ejs.render(templ,filters);
            $('#syslog_filter_query tr').remove();
            $('#syslog_filter_query').append(html);

            //全选
            $("#filters_all_checker").on("click", function () {
                var all_checked =  this.checked;
                $(".filters_row_checker").each(function(){
                    this.checked = all_checked;
                    return true;
                });
                return true;
            });
            //批量删除
            $("#filters_delete").on("click", function () {
                if (!$(".filters_row_checker:checked").length){
                    bootbox.alert('请选择一条记录！');
                    return
                }
                bootbox.confirm("确认删除选定信息？", function(result){
                    if (!result) {
                        return
                    }
                    $(".filters_row_checker:checked").each(function () {
                        deleteFilterById(this);
                    });
                });
                return false
            });
            $('#filters_edit').on('click',function () {
                var elements = $(".filters_row_checker:checked");
                if (elements.length == 1) {
                    var id = elements.first().attr("checked_filter_id");
                    window.location.href = "edit.html?id="+id;
                } else if (elements.length == 0) {
                    bootbox.alert('请选择一条记录！');
                } else {
                    bootbox.alert('你选择了多条记录，请选择一条记录！');
                }
                return false
            });

        },
        error:function (XmlHttpRequest) {
            showErrors(XmlHttpRequest);
            return;
        }
    });
});
//删除
function deleteFilterById(abutt) {
    var id = $(abutt).attr("filter_id");
    var param = JSON.stringify(id);
    $.ajax({
        url:"/filters/"+id,
        type:"DELETE",
        data:param,
        dataType : "json",
        contentType: "application/json",
        error:function (XmlHttpRequest) {
            if (XmlHttpRequest.responseText=='OK'){
                location.reload(true);
            }else {
                showErrors(XmlHttpRequest);
                return;
            }
        }
    });
}

//显示来自服务端的错误
function showErrors(err) {
    var templ = document.getElementById('errors_templ').innerHTML;
    var data = {errors:err };
    var html = ejs.render(templ, data);
    $('#syslog_error_brother+div').remove();
    $('#syslog_error_brother').after(html);
}

//将错误信息显示在页面上 需要传入一个字符串数组
function showMessages(messages) {
    var templ = document.getElementById('filter_error_messages').innerHTML;
    var errors = {messages:messages};
    var html = ejs.render(templ,errors);
    $('#syslog_error_brother+div').remove();
    $('#syslog_error_brother').after(html);
}

