//----------------------------------------------------------------------
//  モード
//----------------------------------------------------------------------
"use strict";

//----------------------------------------------------------------------
//  モジュール読み込み
//----------------------------------------------------------------------
const gulp = require("gulp");
const { src, dest, series, parallel, watch, tree } = require("gulp");
const sass = require("gulp-sass")(require("sass"));

const bs = require("browser-sync");

//----------------------------------------------------------------------
//  関数定義
//----------------------------------------------------------------------
function bsInit(done) {
  bs.init({
    proxy: "localhost:8050",
    notify: false,
    open: false,
    injectChanges: true,
    files: [
      "./app/templates/**/*.html",
      "./app/css/**/*.css",
      "./app/js/**/*.js",
    ],
    port: 3000,
    ui: {
      port: 3001,
    },
  });

  done();
}

function bsReload(done) {
  bs.reload();
  done();
}

function compileSass() {
  return src("./app/scss/**/*.scss")
    .pipe(sass().on("error", sass.logError))
    .pipe(dest("./app/css"));
}

function watchTask(done) {
  // SCSSファイルの監視とコンパイル
  watch("./app/scss/**/*.scss", series(compileSass, bsReload));
  // HTMLファイルの監視
  watch("./app/templates/**/*.html", series(bsReload));
  // JSファイルの監視
  watch("./app/js/**/*.js", series(bsReload));
}

//----------------------------------------------------------------------
//  タスク定義
//----------------------------------------------------------------------
exports.bs = series(compileSass, bsInit, watchTask);

/************************************************************************/
/*  END OF FILE                                                         */
/************************************************************************/
