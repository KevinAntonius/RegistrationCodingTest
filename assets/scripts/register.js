$(document).ready(function(){
    $("#loginbtn").click(function(e) {
        e.preventDefault();
        var email = $("#email").val();
        $.ajax({
            url: '/api/login',
            type: "GET",
            dataType: 'json',
            data: {"email":email},
            success: function(data)
            {
                var name = data.fname;
                window.open("http://localhost:8080/login?fname="+ name,"_self");
            },
            error: function (xhr, err) {
                var msg = xhr.responseText
                $(".message").text(msg)
                $(".alert").show();
            }
        });
    });
    $(".closebtn").click(function(e){
        $(".alert").hide();
    });
    $("form").submit(function(e) {
        e.preventDefault();
        
        var form = $(this);
        var url = form.attr('action');

        $.ajax({
            url: '/api/register',
            type: "POST",
            contentType: 'application/json',
            data: formToJSON(form),
            beforeSend: function(){
                grayout(true);
            },
            success: function()
            {
                $(".logincontainer").show();
            },
            error: function (xhr, err) {
                var msg = ""
                if(xhr.status == '409'){
                    msg = translateError(xhr.status, xhr.responseText)
                } else{
                    msg = xhr.responseText
                }
                $(".message").text(msg)
                $(".alert").show();
                grayout(false);
            }
        });
    });


});

function formToJSON(data){
    var temp_array = {};
    var form_array = data.serializeArray();
    $.map(form_array, function(n, i){
        if (n['value'] != ''){
            temp_array[n['name']] = n['value'];
        }
    });
    return JSON.stringify(temp_array)
}

function translateError(status, responseText){
    var temp = '';
    var arr = responseText.split(" ");
    var field = '';
    arr[1] = arr[1].slice(0,arr[1].length-1)
    console.log(arr)
    if (arr[0] == 'null'){
        temp = 'field ### is empty'
    } else if (arr[0] == 'duplicate'){
        temp = '### is already registered'
    } else if (arr[0] == 'not-valid'){
        if (arr[1] == 'mobile'){
            temp = 'please enter valid Indonesian ###'
        } else{
            temp = 'please enter valid ### address'
        }
        
    }
    if (arr[1] == 'mobile'){
        field = 'phone number'
    } else if (arr[1] == 'fname'){
        field = 'first name'
    } else if (arr[1] == 'lname'){
        field = 'last name'
    } else if (arr[1] == 'email'){
        field = 'email'
    }
    var result = temp.replace('###',field)
    return result.charAt(0).toUpperCase() + result.slice(1)
}

function grayout(state){
    $("input").attr("disabled", state);
    $("#registerbtn").attr("disabled", state);
    if (state == true){
        $(".regcontainer").css({ 
            background: "lightgray",
            opacity: 0.5,
        });
    } else{
        $(".regcontainer").css({ 
            background: '',
            opacity: '',
        });
    }
}