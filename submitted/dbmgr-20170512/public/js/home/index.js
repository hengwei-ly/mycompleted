(function ($) {

  $(function() {
    var wsData = {
      labels: ["财务", "行政", "运维", "***"],
      datasets: [
        {
          label: "计数",
          backgroundColor: [
            'rgba(255, 99, 132, 0.2)',
            'rgba(54, 162, 235, 0.2)',
            'rgba(255, 206, 86, 0.2)',
            'rgba(75, 192, 192, 0.2)',
            'rgba(153, 102, 255, 0.2)',
            'rgba(255, 159, 64, 0.2)'
          ],
          borderColor: [
            'rgba(255,99,132,1)',
            'rgba(54, 162, 235, 1)',
            'rgba(255, 206, 86, 1)',
            'rgba(75, 192, 192, 1)',
            'rgba(153, 102, 255, 1)',
            'rgba(255, 159, 64, 1)'
          ],
          data: [65, 59, 80, 81]
        }
      ]
    };

    var ctx2 = document.getElementById("wsChart").getContext("2d");
    new Chart(ctx2, {type: 'bar', data: wsData, options:{responsive: true, legend: {
      hidden:true,
      display:false
    } }});


    var swData = {
      labels: ["数据库", "中间件", "财务", "OA", "***"],
      datasets: [
        {
          label: "计数",
          backgroundColor: [
            'rgba(255, 99, 132, 0.2)',
            'rgba(54, 162, 235, 0.2)',
            'rgba(255, 206, 86, 0.2)',
            'rgba(75, 192, 192, 0.2)',
            'rgba(153, 102, 255, 0.2)',
            'rgba(255, 159, 64, 0.2)'
          ],
          borderColor: [
            'rgba(255,99,132,1)',
            'rgba(54, 162, 235, 1)',
            'rgba(255, 206, 86, 1)',
            'rgba(75, 192, 192, 1)',
            'rgba(153, 102, 255, 1)',
            'rgba(255, 159, 64, 1)'
          ],
          data: [65, 32, 80, 81, 33]
        }
      ]
    };

    ctx2 = document.getElementById("swChart").getContext("2d");
    new Chart(ctx2, {type: 'bar', data: swData, options:{responsive: true, legend: {
      hidden:true,
      display:false,
    } }});

    var contractData = {
      labels: ["7天内到期", "30天内到期", "3个月内到期", "***"],

      datasets: [
        {
          label: "计数",
          backgroundColor: [
            'rgba(255, 99, 132, 0.2)',
            'rgba(54, 162, 235, 0.2)',
            'rgba(255, 206, 86, 0.2)',
            'rgba(75, 192, 192, 0.2)',
            'rgba(153, 102, 255, 0.2)',
            'rgba(255, 159, 64, 0.2)'
          ],
          borderColor: [
            'rgba(255,99,132,1)',
            'rgba(54, 162, 235, 1)',
            'rgba(255, 206, 86, 1)',
            'rgba(75, 192, 192, 1)',
            'rgba(153, 102, 255, 1)',
            'rgba(255, 159, 64, 1)'
          ],
          data: [65, 32, 80, 81, 33]
        }
      ]
    };

    ctx2 = document.getElementById("contractChart").getContext("2d");
    new Chart(ctx2, {type: 'bar', data: contractData, options:{responsive: true, legend: {
      hidden:true,
      display:false
    }}});
  })

})(jQuery);