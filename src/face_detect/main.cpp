#include <opencv2/opencv.hpp>
#include <filesystem>
#include <iostream>
#include <random>
#include "cv_algorithm.hpp"

constexpr int B_MIN = 80;
constexpr int B_MAX = 200;
constexpr int G_MIN = 80;
constexpr int G_MAX = 220;
constexpr int R_MIN = 160;
constexpr int R_MAX = 255;
constexpr int H_MAX = 30;
constexpr int H_MIN = 0;
constexpr int S_MAX = 150;
constexpr int S_MIN = 30;
constexpr int V_MAX = 255;
constexpr int V_MIN = 60;

#define MIN_HSVCOLOR cv::Scalar(0, 60, 80)
#define MAX_HSVCOLOR cv::Scalar(10, 160, 240)

#define MIN_BKCOLOR cv::Scalar(0, 0, 0)
#define MAX_BKCOLOR cv::Scalar(180, 255, 40)

//1箱の限界値を定める
constexpr double img_width = 730;
constexpr double img_height = 730;
constexpr double resize_param = 1;
constexpr double resize_center = resize_param * 2;
constexpr int offset_width = 180;
constexpr int offset_height = 180;

cv::Mat PinP_point(const cv::UMat &srcImg, const cv::UMat &smallImg, const cv::Point2f p0, const cv::Point2f p1)
{

	cv::Mat dstImg;
	srcImg.copyTo(dstImg);

	std::vector<cv::Point2f> src, dst;
	src.push_back(cv::Point2f(0, 0));
	src.push_back(cv::Point2f(smallImg.cols, 0));
	src.push_back(cv::Point2f(smallImg.cols, smallImg.rows));

	dst.push_back(p0);
	dst.push_back(cv::Point2f(p1.x, p0.y));
	dst.push_back(p1);

	cv::Mat mat = cv::getAffineTransform(src, dst);

	cv::warpAffine(smallImg, dstImg, mat, dstImg.size(), cv::INTER_LINEAR, cv::BORDER_TRANSPARENT);
	return dstImg;
}

int main(int argc, char *argv[])
{
	std::string filePath;
	if (argc < 2)
	{
		std::cout << "err" << std::endl;
		filePath = "/template_1.jpg";
	}
	else
	{
		filePath = argv[1];
	}
	std::string filename;
	cv::UMat result_img;

	cv::Mat template_img = cv::imread(filePath);
	template_img.copyTo(result_img);
	//cv::Mat test_img = cv::imread("images/test_img.png");
	int width = 100;
	int height = 100;
	int index_w = 0;
	int index_h = 0;
	int size_w = template_img.cols;
	int size_h = template_img.rows;

	std::random_device rnd; // 非決定的な乱数生成器でシード生成機を生成
	std::mt19937 mt(rnd()); //  メルセンヌツイスターの32ビット版、引数は初期シード
	//ここで画像の配置候補を作成する

	int cnt_h = static_cast<int>(template_img.rows / img_height);
	int cnt_w = static_cast<int>(template_img.cols / img_width);

	std::uniform_int_distribution<> rand_w(0, cnt_w - 2); // [0, 99] 範囲の一様乱数
	std::uniform_int_distribution<> rand_h(0, cnt_h - 2); // [0, 99] 範囲の一様乱数
	std::vector<std::vector<bool>> img_map(cnt_w, std::vector<bool>(cnt_h));

	for (const std::filesystem::directory_entry &i : std::filesystem::directory_iterator("/cheese/images"))
	{

		filename = "/cheese/images/" + i.path().filename().string();

		std::cout << i.path().filename().string() << std::endl;
		cv::UMat _extract_img;
		cv::UMat extract_img;

		extractFaceImage(filename, _extract_img);
		if (_extract_img.empty())
		{
			std::cout << "skipped" << std::endl;
			continue;
		}
		//リサイズをかける
		//小倉駅のみスケールする
		double scale = filePath == "/template_3.jpg" ? 0.7 : 1.0;
		cv::resize(_extract_img, extract_img, cv::Size(), (0.85 * scale * img_width) / _extract_img.cols, (0.85 * scale * img_width) / _extract_img.cols);

		//配置場所を決める

		do
		{
			index_w = rand_w(rnd);
			index_h = rand_h(rnd);
		} while (img_map.at(index_w).at(index_h));

		img_map.at(index_w).at(index_h) = false;

		cv::circle(result_img, cv::Point2f(img_width * index_w + (extract_img.cols / 2.0) + offset_width, img_height * index_h + (extract_img.rows / 2.0) + offset_height), 450 * (scale + 0.05), cv::Scalar(240, 240, 240), -1);
		auto tmp = PinP_point(result_img, extract_img, cv::Point2f(img_width * index_w + offset_width, img_height * index_h + offset_height), cv::Point2f(img_width * index_w + extract_img.cols + offset_width, img_height * index_h + extract_img.rows + offset_height));
		tmp.copyTo(result_img);
	}
	cv::imwrite("/cheese/result.jpg", result_img);
	//プレビュー用軽量化画像の生成
	cv::UMat mini_result;
	cv::resize(result_img, mini_result, cv::Size(), 0.5, 0.5);
	cv::imwrite("/cheese/result_mini.jpg", mini_result);

	std::cout << "success make file" << std::endl;
}
