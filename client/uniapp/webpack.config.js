const path = require('path');

module.exports = {
    mode: 'development',

    entry: './client.ts',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist')
    },

    module: {
        rules: [{
            test: /\.ts$/,
            use: "ts-loader"
        }]
    },
    resolve: {
        extensions: [
            '.ts',
            '.js'
        ]
    }
};