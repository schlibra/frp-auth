<?php

namespace app\controller;

use app\BaseController;
use app\Request;
use think\facade\Log;
use think\response\Json;
use WpOrg\Requests\Requests;

class Index extends BaseController
{
    public function handler(Request $request): Json {
        $op = $request->param('op');
        if ($op === "Login") {
            $content = $request->post('content');
            $user = $content['user'];
            if (isset($content['metas']) && isset($content['metas']['secret'])) {
                $secret = $content['metas']['secret'];
                $res = Requests::post('https://frp.schhz.cn/api/auth', [], [
                    'username' => $user,
                    'secret' => $secret,
                ])->body;
                Log::write($res);
                $data = json_decode($res);
                if ($data->code === 200) {
                    return json(['reject' => false, "unchange" => true]);
                } else {
                    return json(['reject' => true, "reject_reason" => "用户名和秘钥不匹配"]);
                }
            } else {
                return json(['reject' => true, "reject_reason" => "未填写秘钥"]);
            }
        } else {
            return json(['reject' => true]);
        }
    }
}
