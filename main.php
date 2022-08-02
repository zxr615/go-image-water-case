<?php

$t1 = microtime(true);

$dst  = __DIR__ . '/test.png';
$src  = __DIR__ . '/test_water.png';
$font = __DIR__ . '/font.ttf';

$t3    = microtime(true);
$srcIm = imagecreatefrompng($src);
$t4    = microtime(true);

$num = 100;
for ($i = 0; $i < $num; $i++) {
    $t5    = microtime(true);
    $dstIm = imagecreatefrompng($dst);
    $t6    = microtime(true);

    $dstSize = getimagesize($dst);
    $srcSize = getimagesize($src);
    $new     = imagecreatetruecolor($dstSize[0], $dstSize[1]);
    imagecopy($new, $dstIm, 0, 0, 0, 0, $dstSize[0], $dstSize[1]);
    imagecopy($new, $srcIm, $dstSize[0] - 200, $dstSize[1] - 220, 0, 0, $srcSize[0], $srcSize[1]);
    $rgb = imagecolorallocate($dstIm, 21, 33, 57);//字体颜色
    $id  = $i + 1000;
    imagefttext($new, 30, 0, $dstSize[0] - 200, $dstSize[1] - 20, $rgb, $font, "ID: {$id}");
    $t7 = microtime(true);
    imagejpeg($new, __DIR__ . "/water/{$i}.jpeg");
    $t8 = microtime(true);
    imagedestroy($dstIm);
    imagedestroy($new);
}

imagedestroy($srcIm);
$t2 = microtime(true);

echo '总耗时: ' . round($t2 - $t1, 5) . 's' . PHP_EOL;
echo '总耗时: ' . round($t2 - $t1, 5) * 1000 . 'ms' . PHP_EOL;
echo '水印decode耗时: ' . round($t4 - $t3, 5) * 1000 . 'ms' . PHP_EOL;
echo '原图decode耗时: ' . round($t6 - $t5, 5) * 1000 . 'ms' . PHP_EOL;
echo '新图encode耗时: ' . round($t8 - $t7, 5) * 1000 . 'ms' . PHP_EOL;
