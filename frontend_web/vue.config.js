const path = require('path')

function resolve(dir) {
  return path.join(__dirname, dir)
}

module.exports = {
  publicPath: '/',
  outputDir: 'dist',
  assetsDir: 'static',
  productionSourceMap: false,
  configureWebpack: {
    name: 'YUTANK Shop',
    resolve: {
      alias: {
        '@': resolve('src')
      }
    }
  },
  devServer: {
    port: 9528,
    open: false,
    overlay: {
      warnings: false,
      errors: true
    }
  }
}
