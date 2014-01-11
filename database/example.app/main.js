var app = {
    Run: _.once(function() {
        var isLoaded = false;
        _.delay(function() {
            isLoaded = true;
        }, 4000);

        var loadingFn = function() {
            if (isLoaded) {
                $("body pre").text("Loaded");
                return;
            }

            $("body pre").append(".");

            _.delay(loadingFn, 500);
        };
        loadingFn();
    })
};

app.Run();
