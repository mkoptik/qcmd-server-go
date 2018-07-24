// development config
const merge = require('webpack-merge');
const webpack = require('webpack');
const commonConfig = require('./common');

module.exports = merge(commonConfig, {
    mode: 'development',
    entry: {
        client: './index.tsx'
    },
    devtool: 'cheap-module-eval-source-map',
    plugins: [
        new webpack.NamedModulesPlugin(), // prints more readable module names in the browser console on HMR updates
    ],
});