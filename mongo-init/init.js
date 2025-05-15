db.createCollection("measurements", {
    timeseries: {
        timeField: "timestamp",
        metaField: "asset_id",
        granularity: "seconds"
    }
});
