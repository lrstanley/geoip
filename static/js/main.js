$(".lookup").click(function(event){
    window.location.href = "/" + $("#query").val();
});

$('#query').keypress(function (e) {
  if (e.which == 13) {
    window.location.href = "/" + $("#query").val();
  }
});