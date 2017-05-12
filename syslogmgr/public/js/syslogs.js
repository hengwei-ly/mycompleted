/**
 * Created by pc on 2017/3/29.
 */
$(function(){
    // get the EJS template as a string
    var templ = document.getElementById('body').innerHTML;
    // data to output to the template function
    var size = $("#syslog_pagesize_select").val();
    console.log(size);
    $.get("hengwei/aaa/Syslogs/GetCounts",function (data) {
        var count = eval(data);
        var page = 1;

        var maxpage = Math.ceil(count/size);
        $.get("hengwei/aaa/Syslogs/QueryJSON",function (datas) {
            var data = { names: datas };
            var html = ejs.compile(templ)(data);
            document.body.innerHTML = html;
        });
    });
});