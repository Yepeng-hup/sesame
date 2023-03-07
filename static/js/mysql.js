function pushFrom() {
    var my_user=document.getElementById('my-user').value;
    var my_passwd=document.getElementById('my-passwd').value;
    var my_dbname=document.getElementById('my-dbname').value;
    var my_dbip=document.getElementById('my-dbip').value;
    var my_dbport=document.getElementById('my-dbport').value;
    var my_sql=document.getElementById('my-sql').value;
    var my_bool=document.getElementById('my-bool').value;
    var my_threadingNum=document.getElementById('my-threadingNum').value;
    var my_Qnum=document.getElementById('my-Qnum').value;
    var my_sleepNum=document.getElementById('my-sleepNum').value;

    switch (true) {
        case my_sql.includes("insert"):
            break;
        case my_sql.includes("INSERT"):
            break;
        default:
            window.location = '/index';
            alert("Error Only sql is supported 'insert'!!!");
            return
    }

    if (my_passwd == ""||my_dbname == ""||my_sql == ""){
        window.location = '/index';
        alert("Field cannot be empty！");
    }else {

        if (my_bool == "true") {
            $.ajax({
                url: "/mysql/use",
                type: "POST",
                data: {
                    "user": my_user,
                    "passwd": my_passwd,
                    "dbName": my_dbname,
                    "dbIp": my_dbip,
                    "dbPort": my_dbport,
                    "sql": my_sql,
                    "bools": my_bool,
                    "threadingNum": my_threadingNum,
                    "Qnum": my_Qnum,
                    "sleepNum": my_sleepNum
                }
            });
            window.location = "/index";
            alert("success,at work!");

        } else {
            $.ajax({
                url: "/mysql/use",
                type: "POST",
                data: {
                    "user": my_user,
                    "passwd": my_passwd,
                    "dbName": my_dbname,
                    "dbIp": my_dbip,
                    "dbPort": my_dbport,
                    "sql": my_sql,
                    "bools": my_bool,
                    "Qnum": my_Qnum,
                    "sleepNum": my_sleepNum
                },
            });
            window.location = "/index";
            alert("success,at work!");
        }
    }
}