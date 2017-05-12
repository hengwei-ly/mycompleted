/**
 * Created with Gogland.
 * User: eeiceyes
 * Date: 2017/2/27
 * Time: 9:08
 * To change this template use File | Settings | File Templates.
 */

(function ($) {
  window.App = {
    deleteWith : function(url, callback, descr) {
      var self = this;

      var callFn = function(){};
      var opDescr = "操作失败";
      if(typeof callback === "function") {
        callFn = callback;
      }
      if(typeof descr === "string") {
        opDescr = "操作[" + descr + "]失败\r\n";
      } else if(typeof callback === "string") {
        opDescr = "操作[" + callback + "]失败\r\n";
      }

      $.ajax({
        url: url,
        method: "delete",
        success: callFn,
        error: function(error) {
          Msg.error(opDescr + self.resolveAjaxError(error));
        }
      })
    },
    linkTo : function(url) {
      var newUrl = url.indexOf("/") == 0 ? url : ("/" + url);
      var fullUrl = url.indexOf("http") == 0 ? url : (window.location.origin + newUrl);

      window.location.href = fullUrl;
    },
    resolveAjaxError: function(error) {
      return error.status + "\r\n" + error.statusText + "\r\n" + error.responseText;
    },
    serializeJson: function(objArray){
      return JSON.stringify(this.serializeObject(objArray));
    },
    serializeObject: function(objArray, movePrefix){
      var value = {};
      for(var i=0; i<objArray.length; i++) {
        var v = objArray[i];
        var name = v.name;

        if(movePrefix && name.indexOf(movePrefix) == 0)
          name = name.replace(movePrefix, "");

        if(!value[name]) {
          value[name] = [];
        }
        if(name.indexOf("[]") >= 0)
          value[name].push(v.value);
        else
          value[name] = v.value;
      }
      return value;
    }
  };

  window.Msg = {
    success : function(msg) {
      toastr.success(msg);
    },
    info : function(msg) {
      toastr.info(msg);
    },
    warn : function(msg) {
      toastr.warning(msg);
    },
    error : function(msg) {
      toastr.error(msg);
    }
  };

  window.Modal = {
    show : function(selector, callback, options) {
      var innerForm = $(selector).find("form");
      if(options.resetForm && innerForm.size() > 0 && typeof innerForm[0].reset === 'function') {
        innerForm[0].reset();
      }

      if(typeof options.data !== 'undefined') {
        this._fillForm(innerForm, options.data)
      }

      $(document.body).undelegate("#btnOk", "click");
      $(document.body).delegate("#btnOk", "click", function() {
        if(innerForm.valid()) {
          var formData = innerForm.serialize();
          var dataArray = innerForm.serializeArray();
          var jsonObject = App.serializeObject(dataArray);
          var jsonData = JSON.stringify(jsonObject);
          callback({form:formData, array:dataArray, json:jsonData, jsonObject:jsonObject});

          $(selector).find("[data-dismiss='modal']").click();
        }
      });

      $(selector).one('shown.bs.modal', function () {
        //to be fixed
      });
      $(selector).one('hide.bs.modal', function () {
        //to be fixed
      });

      $(selector).modal({backdrop:'static'});
    },
    _fillForm: function(form, data) {
      for(var key in data) {
        var selector = "[name='" + key + "']";
        form.find(selector).val(data[key]);
      }
    },
    _validateForm : function(form) {
      var custom = {
        focusCleanup: false,

        wrapper: 'div',
        errorElement: 'span',

        highlight: function (element) {
          $(element).parents('.form-group').removeClass('success').addClass('error');
        },
        success: function (element) {
          $(element).parents('.form-group').removeClass('error').addClass('success');
          var wrap = $(element).parents('.controls:not(:has(.clean))');
          if(wrap.find(".help-block").size() > 0)
            wrap.find('.help-block').before('<div class="clean"></div>');
          else
            wrap.find('div:last').before('<div class="clean"></div>');
        },
        errorPlacement: function (error, element) {
          if (element.parents('.controls').size() > 0)
            error.appendTo(element.parents('.controls'));
          else
            error.appendTo(element.parent());
        }
      };
      if($.validator) {
        $.validator.autoCreateRanges = true;
        $(form).validate(custom);
      }
    }
  };
})(jQuery);

(function ($) {
  $(document.body).delegate("a[data-toggle='submit']", "click", function(event) {
    event.preventDefault();

    var targetSelector = $(this).attr("data-target");
    $(targetSelector).attr("action", $(this).attr("href"));

    var method = $(this).attr("method");
    $(targetSelector).attr("method", "POST");
    if(typeof method === 'undefined' || method.toLowerCase() == "get") {
      $(targetSelector).attr("method", "GET");
    } else {
      $(targetSelector).find("[name=_method]").val(method);
    }

    $(targetSelector).submit();
  })
})(jQuery);