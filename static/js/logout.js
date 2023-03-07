function logout() {
    $.ajax({
        url: "/logout",
        type: "GET",
        success: function(data) {
            if (data.status == 200) {
                alert("退出登录成功!");
                window.location = '/';
            } else {
                alert("退出登录失败!");
            }
        },
        error: function() {
            alert("请求失败！");
        }
    });
}