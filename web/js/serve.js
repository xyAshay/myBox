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
            $.ajax({
                url:"/api/delete/"+path,
                type: "DELETE",
                success(){
                    alert("File " + path +" deleted");
                    location.reload();
                }
            });
        }
    }
})

