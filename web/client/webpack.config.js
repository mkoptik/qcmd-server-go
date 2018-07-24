const Webpack = require('webpack');
const Path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const isProduction = process.argv.indexOf('-p') >= 0;
const outPath = Path.join(__dirname, './dist/');
const sourcePath = Path.join(__dirname, './src');

module.exports = {
    context: sourcePath,
    entry: {
        client: './index.tsx',
        vendor: [
            'react',
            'react-dom',
            'react-router',
            'react-router-dom',
            'react-autosuggest',
            'axios'
        ]
    },
    output: {
        path: outPath,
        publicPath: '/',
        filename: '[name].js',
        libraryTarget: 'var',
        library: 'QcmdClientEntryPoint'
    },
    target: 'web',
    resolve: {
        extensions: ['.js', '.ts', '.tsx'],
        // Fix webpack's default behavior to not load packages with jsnext:main module
        // https://github.com/Microsoft/TypeScript/issues/11677
        mainFields: ['browser', 'main']
    },
    module: {
        loaders: [
            // .ts, .tsx
            {
                test: /\.tsx?$/,
                use: isProduction
                    ? ['awesome-typescript-loader?module=es5']
                    : ['react-hot-loader/webpack', 'awesome-typescript-loader']
            },
            // css
            {
                test: /\.css$/,
                use: ExtractTextPlugin.extract({
                    fallback: 'style-loader',
                    use: [
                        {
                            loader: 'css-loader',
                            query: {
                                url: false,
                                modules: true,
                                sourceMap: !isProduction,
                                importLoaders: 1,
                                localIdentName: '[local]'
                            }
                        },
                        {
                            loader: 'postcss-loader',
                            options: {
                                ident: 'postcss',
                                plugins: [
                                    require('postcss-import')({addDependencyTo: Webpack}),
                                    require('postcss-url')({url: "inline"}),
                                    require('postcss-cssnext')(),
                                    require('postcss-reporter')(),
                                    require('postcss-browser-reporter')({disabled: isProduction}),
                                ]
                            }
                        }
                    ]
                })
            },
            // static assets
            {test: /\.html$/, use: 'html-loader'},
            {
                test: /\.(png|jpg|gif)$/,
                use: {
                    loader: 'file-loader',
                    options: {
                        name: 'files/[name].[ext]'
                    }
                }
            }
        ]
    },
    plugins: [
        new Webpack.DefinePlugin({
            'process.env.NODE_ENV': isProduction === true ? JSON.stringify('production') : JSON.stringify('development')
        }),
        new Webpack.optimize.CommonsChunkPlugin({
            name: 'vendor',
            filename: 'vendor.bundle.js',
            minChunks: Infinity
        }),
        new Webpack.optimize.AggressiveMergingPlugin(),
        new ExtractTextPlugin({
            filename: 'styles.css',
            disable: !isProduction
        })
    ],
    devServer: {
        contentBase: sourcePath,
        hot: true,
        stats: {
            warnings: false
        }
    },
    node: {
        // workaround for webpack-dev-server issue
        // https://github.com/webpack/webpack-dev-server/issues/60#issuecomment-103411179
        fs: 'empty',
        net: 'empty'
    }
};
