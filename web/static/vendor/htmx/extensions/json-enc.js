htmx.defineExtension('json-enc', {
    onEvent: function (name, evt) {
        if (name === "htmx:configRequest") {
            evt.detail.headers['Content-Type'] = "application/json";
        }
    },

    encodeParameters: function (xhr, parameters, elt) {
        xhr.overrideMimeType('text/json');

        const transformValue = (value, inputElement) => {
            if (inputElement && inputElement.type === "number") {
                return Number(value);
            }
            if (inputElement && inputElement.type === "select-one") {
                return Number(value);
            }
            if (inputElement && inputElement.type === "datetime-local" && typeof value === "string" && value.includes("T")) {
                // Ensure datetime is in the correct format (e.g., add :00Z if needed)
                return value.endsWith("Z") ? value : value + ":00Z";
            }
            return value;
        };

        const transformedParameters = Object.fromEntries(
            Object.entries(parameters).map(([key, value]) => {
                const inputElement = elt.closest('form').querySelector(`[name="${key}"]`);
                return [key, transformValue(value, inputElement)];
            })
        );

        return JSON.stringify(transformedParameters);
    }
});