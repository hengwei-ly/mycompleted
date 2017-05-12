/**
 * Created by pc on 2017/3/27.
 */

$(function () {
    var page = "alerts";
    deleteByPage(page);
    $("#tab_alerts").on("click",function () {
        page = "alerts";
        deleteByPage(page);
    });
    $("#tab_traps").on("click",function () {
        page = "traps";
        deleteByPage(page);
    });
    $("#tab_histories").on("click",function () {
        page = "histories";
        deleteByPage(page);
    });
    $("#tab_syslogs").on("click",function () {
        page = "syslogs";
        deleteByPage(page);
    });
});

function deleteByPage(page) {
    var urlPrefix = $("#urlPrefix").val();

    var checkerAllId = "dev_tables-all-checker-" + page;
    var checkerRowId = "dev_tables-row-checker-" + page;

    $("#"+checkerAllId).on("click", function () {
        var all_checked =  this.checked;

        $("."+checkerRowId).each(function(){
            this.checked = all_checked;
            return true;
        });
        return true;
    });
    $("#dev_tables-delete").on("click", function () {

        bootbox.confirm("确认删除选定信息？", function (result) {
            if (!result) {
                return;
            }

            var f = document.createElement("form");
            f.action = $("#dev_tables-delete").attr("url");
            f.method = "POST";
            var inputField = document.createElement("input");
            inputField.type = "hidden";
            inputField.name = "_method";
            inputField.value = "DELETE";
            f.appendChild(inputField);

            $("." + checkerRowId + ":checked").each(function () {
                var inputField = document.createElement("input");
                inputField.type = "hidden";
                inputField.name = "names[]";
                inputField.value = $(this).attr("key");
                f.appendChild(inputField);
            });

            var inputPage = document.createElement("input");
            inputPage.type = "hidden";
            inputPage.name = "page";
            inputPage.value = page;
            f.appendChild(inputPage);

            document.body.appendChild(f);
            f.submit();
        });
        return false
    });
}

