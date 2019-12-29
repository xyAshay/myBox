new Vue({
    el:"#app",
    data() {
        return{
            files:[{
                name:"loading..",
                size:"loading...",
                path:"loading..."
             }]
        }
    },
    created(){
        $.getJSON(("/api/json"+location.pathname).replace('/serve/','/'), (res) => {
            res.sort((a,b) => (a.type > b.type) ? 1 : ((b.type > a.type) ? -1 : 0)); 
            this.files = res;
        });
    },
    methods:{
        deleteFile(path){
            var r = confirm("Delete "+path+" ??");
            if(r == true){
                $.ajax({
                    url:"/api/delete/"+path,
                    type: "DELETE",
                    success(){
                        alert(path +" Deleted");
                        location.reload();
                    }
                });
            }
            else{
                alert("Deletion Cancelled!!");
                location.reload();
            }
        }
    }
})

$(document).ready(() => {
    $("#confirmUpload").hide();
    $("#newdirform").hide();

    $("#uploader").change(() => {
        $("#uploader").hide();
        $("#confirmUpload").show();
    });

    $("#confirmUpload").click(function (e){
        e.preventDefault();
        $.ajax({
            url: ('/api/upload'+location.pathname).replace('/serve/','/'),
            type: 'POST',
            enctype: 'multipart/form-data',
            contentType: false,
            processData: false,
            data: new FormData($('#uploadform')[0]),
            success: () => {
                alert("File Uploaded");
                location.reload();
            }
        });
    });

    $("#newdir").click(() => {
        $("#newdir").hide();
        $("#newdirform").show();
    });

    $("#create").submit(() => {
        if($("#dirname").val() != ""){
            var targetURI;
            if(location.pathname != '/serve/')
                targetURI = ('/api/create'+location.pathname+'/'+$("#dirname").val()).replace('/serve/','/');
            else
                targetURI = ('/api/create'+'/'+$("#dirname").val());
            $.ajax({
                url: targetURI,
                type: "POST",
                success: () => {
                    alert($("#dirname").val() + " Created");
                    location.reload();
                } 
            });
        }
        else
            alert("Please Enter Folder Name");
    });
});

