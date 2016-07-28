
$(document).ready(function(){


$("#bbb").click(function(){ 
	if ($("#bbb").attr("checked")=="checked") { 
		$(" #ols").show();
		$(".li5").height(430);
	}
	else{
		$("#ols").hide();
		$(".li5").height(60);
	}
}) 

$("#ccc").click(function(){ 
	if ($("#ccc").attr("checked")=="checked") { 
		$(" #ols1").show();
		$(".li6").height(200);
	}
	else{
		$("#ols1").hide();
		$(".li6").height(60);
	}
}) 







$("#ddd").click(function(){ 
	if ($("#ddd").attr("checked")=="checked") { 
		$(" #ols2").show();
		$(".li7").height(150);
	}
	else{
		$("#ols2").hide();
		$(".li7").height(60);
	}
})







$("#eee").click(function(){ 
	if ($("#eee").attr("checked")=="checked") { 
		$(" #ols3").show();
		$(".li8").height(150);
	}
	else{
		$("#ols3").hide();
		$(".li8").height(60);
	}
})











$(".li9").click(function(){
	$(".bigs").toggle();
	if($(".li9-b").hasClass("icon-angle-down")){
		$(".li9-b").removeClass("icon-angle-down").addClass("icon-angle-up");
	}else{
		$(".li9-b").removeClass("icon-angle-up").addClass("icon-angle-down");
	}
});




$(".hov1").click(function(){
	$(".hov1").siblings().toggle();
});

$(".hov2").click(function(){
	$(".hov2").siblings().toggle();
});

$(".hov3").click(function(){
	$(".hov3").siblings().toggle();
});


$(".btn-gd").click(function(){
	$(".black").show();
	$(".bomb1").show();
})

$(".act").click(function(){
	$(".black").hide();
	$(".bomb1").hide();
})

$(".tables-ul").click(function(){
	$(this).hide();
})




$(".btn-kr").click(function(){
	$(".black").show();
	$(".bomb2").show();
})

$(".act").click(function(){
	$(".black").hide();
	$(".bomb2").hide();
})

$(".btn-hg").click(function(){
	$(".black").show();
	$(".bomb3").show();
})

$(".act").click(function(){
	$(".black").hide();
	$(".bomb3").hide();
})






$(".left-p1").click(function(){
	$(".wrap-gl").show();
	$(".wrap-yc").hide();
	$(".left-p1").css({"background":"#fff","color":"#414971"});
	$(".left-p2").css({"background":"#414971","color":"#fff"});
})
$(".left-p2").click(function(){
	$(".wrap-gl").hide();
	$(".wrap-yc").show();
	$(".left-p2").css({"background":"#fff","color":"#414971"});
	$(".left-p1").css({"background":"#414971","color":"#fff"});
})

if($(".wrap-gl").show()){
	$(".left-p1").css({"background":"#fff","color":"#414971"});
	$(".left-p2").css({"background":"#414971","color":"#fff"});
}






var w=$('.box-ylt').width(),
	bs=$('.bs').width(),
	startX=0, mosveX=0, lock=false, left=16;

$(".bs").on('mousedown',function(e){
	left=$('.box-ylt').offset().left;  // box-ylt距离左边的距离
	startX=e.pageX;		//鼠标

	console.log(e.pageX)
//	console.log(e.clientX)
//	console.log(e.screenX)			  
	lock=true;
})
$(".box-ylt").on('mousemove',function(e){
	if(lock){
		moveX=e.pageX-left;
		if(moveX>w) moveX=w;
		if(moveX<0) moveX=0;
		$('.bs').css('marginLeft',moveX/bs-16);
		$('.spans').css('width',moveX);
	}
})
$('.box-ylt').on('mouseup',function(e){
	lock=false;
})












});

// 原生