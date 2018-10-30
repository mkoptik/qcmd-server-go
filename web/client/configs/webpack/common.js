// shared config (dev and prod)
const {resolve} = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.jsx'],
    },
    output: {
        filename: 'client.js',
    },
    context: resolve(__dirname, '../../src'),
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: ['awesome-typescript-loader'],
            },
            {
                test: /\.css$/,
                use: ExtractTextPlugin.extract({
                    fallback: 'style-loader',
                    use: 'css-loader'
                })
            },
            {
                test: /\.(jpe?g|png|gif|svg)$/i,
                loaders: [
                    'file-loader?hash=sha512&digest=hex&name=img/[hash].[ext]',
                    'image-webpack-loader?bypassOnDebug&optipng.optimizationLevel=7&gifsicle.interlaced=false',
                ],
            },
        ],
    },
    performance: {
        hints: false,
    },
    plugins: [
        new ExtractTextPlugin("client.min.css", { allChunks: true }),
        new HtmlWebpackPlugin({
            template: resolve(__dirname, '../../src/index.html'),
            favicon: resolve(__dirname, '../../../assets/favicon.ico')
        })
    ]
};