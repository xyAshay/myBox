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
            var r = confirm("Delete File "+path);
            if(r == true){
                $.ajax({
                    url:"/api/delete/"+path,
                    type: "DELETE",
                    success(){
                        alert("File " + path +" deleted");
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
});

