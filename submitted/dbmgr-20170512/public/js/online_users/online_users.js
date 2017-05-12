$(function () {
    var urlPrefix = $("#urlPrefix").val();

    $("#online_users-all-checker").on("click", function () {
        var all_checked =  this.checked;
        $(".online_users-row-checker").each(function(){
            this.checked = all_checked;
            return true;
        });
        return true;
    });


    $("#online_users-delete").on("click", function () {
        bootbox.confirm("确认删除选定信息？", function(result){
            if (!result) {
                return;
            }

            var f = document.createElement("form");
            f.action = $("#online_users-delete").attr("url");
            f.method="POST";
            var inputField = document.createElement("input");
            inputField.type = "hidden";
            inputField.name = "_method";
            inputField.value = "DELETE";

            $(".online_users-row-checker:checked").each(function (i) {
                var inputField = document.createElement("input");
                inputField.type = "hidden";
                inputField.name = "id_list[]";
                inputField.value = $(this).attr("key");
                f.appendChild(inputField);
            });

            document.body.appendChild(f);
            f.submit();
        });
        return false
    });

    $("#online_users-edit").on("click", function () {
        var elements = $(".online_users-row-checker:checked");
        if (elements.length == 1) {
            window.location.href= elements.first().attr("url");
        } else if (elements.length == 0) {
            bootbox.alert('请选择一条记录！')
        } else {
            bootbox.alert('你选择了多条记录，请选择一条记录！')
        }

        return false
    });
});
