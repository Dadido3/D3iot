var api = (function () {
    return {
        getChannels: function () {
            return $.ajax({
                type: "GET",
                url: "/api/getChannels",
                cache: false,
                contentType: 'application/json',
            });
        },

        LAB2sRGB: function (labColor) {
            return $.ajax({
                type: "POST",
                url: "/api/LAB2sRGB",
                cache: false,
                data: JSON.stringify(labColor),
                contentType: 'application/json',
            });
        },

        setDCSVector: function (dcsVector) {
            return $.ajax({
                type: "POST",
                url: "/api/setDCSVector",
                cache: false,
                data: JSON.stringify(dcsVector),
                contentType: 'application/json',
            });
        },

        DCS2LAB: function (dcsVector) {
            return $.ajax({
                type: "POST",
                url: "/api/DCS2LAB",
                cache: false,
                data: JSON.stringify(dcsVector),
                contentType: 'application/json',
            });
        },

        addDataPoint: function (dcsVector, labColor) {
            return $.ajax({
                type: "POST",
                url: "/api/addDataPoint",
                cache: false,
                data: JSON.stringify({LinDCSVector: dcsVector, L: labColor.L, A: labColor.A, B: labColor.B}),
                contentType: 'application/json',
            });
        },

    };
})();