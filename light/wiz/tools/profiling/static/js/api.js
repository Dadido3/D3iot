var api = (function () {
    return {
        setRGBW: function (data) {
            return $.ajax({
                type: "POST",
                url: "/api/setRGBW",
                cache: false,
                data: JSON.stringify(data),
                contentType: 'application/json',
            });
        },

        addDataPoint: function (data) {
            return $.ajax({
                type: "POST",
                url: "/api/addDataPoint",
                cache: false,
                data: JSON.stringify(data),
                contentType: 'application/json',
            });
        },

    };
})();