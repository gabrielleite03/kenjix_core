<?php

require __DIR__ . '/vendor/autoload.php';

use NFePHP\DA\NFe\Danfe;

$xmlPath = $argv[1];
$pdfPath = $argv[2];

if (!file_exists($xmlPath)) {
    echo "XML não encontrado\n";
    exit(1);
}

$xml = file_get_contents($xmlPath);

try {
    $danfe = new Danfe($xml);

    // 🔥 método correto
    $pdf = $danfe->render();

    file_put_contents($pdfPath, $pdf);

    echo "DANFE gerado com sucesso\n";
} catch (Exception $e) {
    echo "Erro: " . $e->getMessage() . "\n";
    exit(1);
}